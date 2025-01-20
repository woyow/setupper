package translate

import (
	"os"
	"unsafe"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type LocalizationMap map[string]map[string]string

type localizationData struct {
	mu   sync.RWMutex
	data LocalizationMap
}

// repo - Need to implement this interface
type repo interface {
	GetLocalizationMap(ctx context.Context, dbName string) (LocalizationMap, error)
}

type Translate struct {
	repo                repo
	localizationData    localizationData
	localizationDataNew localizationData
	updateTimeout       time.Duration
	log                 *logrus.Logger
	dbName              string
	defaultLanguage     string
	availableLanguages  []string
	stopCh              <-chan os.Signal
}

func NewTranslate(repo repo, cfg Config, stopCh <-chan os.Signal, log *logrus.Logger) *Translate {
	t := &Translate{
		repo:               repo,
		localizationData: localizationData{
			mu:   sync.RWMutex{},
			data: make(LocalizationMap, 256),
		},
		localizationDataNew: localizationData{
			mu:   sync.RWMutex{},
			data: make(LocalizationMap, 256),
		},
		updateTimeout:      cfg.UpdateTimeout,
		dbName:             cfg.DBName,
		defaultLanguage:    cfg.DefaultLanguage,
		availableLanguages: cfg.AvailableLanguages,
		log:                log,
		stopCh:             stopCh,
	}

	if err := t.update(); err != nil {
		panic(err)
	}

	go t.updater()

	return t
}

const (
	translateNotExists = "translate_not_exists"
)

func (t *Translate) Translate(key, lang string) string {
	switch {
	case func() bool {
		for _, v := range t.availableLanguages {
			if v == lang {
				return true
			}
		}
		return false
	}():
	default:
		lang = t.defaultLanguage
	}

	t.localizationData.mu.RLock()
	defer t.localizationData.mu.RUnlock()

	val, ok := t.localizationData.data[key][lang]
	if !ok {
		return translateNotExists
	}

	return val
}

func swap(old *LocalizationMap, new *LocalizationMap) {
    p := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(old)), *(*unsafe.Pointer)((unsafe.Pointer(new))))
    _ = atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(new)), p)
}

func (t *Translate) update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	m, err := t.repo.GetLocalizationMap(ctx, t.dbName)
	if err != nil {
		t.log.Error("translate: update - t.repo.GetLocalizationMap error: ", err)
		return err
	}

	t.localizationDataNew.mu.Lock()
	defer t.localizationDataNew.mu.Unlock()
	
	// Clear old keys
	for k := range t.localizationDataNew.data {
		delete(t.localizationDataNew.data, k)
	}

	// Append new keys and values
	for k, v := range m {
		t.localizationDataNew.data[k] = v
	}

	// Swap maps to avoid locking by mutex the main localization map for a long time
	swap(&t.localizationData.data, &t.localizationDataNew.data)

	return nil
}

func (t *Translate) updater() {
	tt := time.NewTicker(t.updateTimeout)
	defer tt.Stop()

	for {
		select {
		case <-tt.C:
			if err := t.update(); err != nil {
				continue
			}
		}
	}
}
