package hashtable

import (
	"crypto/sha256"
	"errors"
)

type keyValuePair struct {
	key      string
	value    int
	nextPair *keyValuePair
	headPair *keyValuePair
	tailPair *keyValuePair
}

type Hashtable struct {
	// TODO: Implement
	table []keyValuePair
}

func New(size int) *Hashtable {
	// TODO: Implement
	h := Hashtable{table: make([]keyValuePair, size)}
	return &h
}

// Insert inserts a new key/value pair into the hashtable.
// Should return an error if the key already exists.
func (h *Hashtable) Insert(key string, value int) error {
	// TODO: Implement
	keyIndex := generateKeyIndex(key)
	emptyPair := keyValuePair{}
	if h.table[keyIndex] == emptyPair {
		h.table[keyIndex] = keyValuePair{key: key, value: value}
	} else {
		newPair := keyValuePair{key: key, value: value}
		if h.table[keyIndex].headPair == nil {
			h.table[keyIndex].headPair = &newPair
			h.table[keyIndex].tailPair = &newPair
		} else {
			h.table[keyIndex].tailPair.nextPair = &newPair
			h.table[keyIndex].tailPair = &newPair
		}
	}
	return nil
}

// Update updates an existing key to be associated with a different value.
// Should return an error if the key doesn't already exist.
func (h *Hashtable) Update(key string, value int) error {
	// TODO: Implement
	keyIndex := generateKeyIndex(key)
	if h.table[keyIndex].key == key {
		h.table[keyIndex].value = value
		for h.table[keyIndex].headPair != nil {
			if h.table[keyIndex].headPair.key == key {
				h.table[keyIndex].headPair.value = value
			}
			h.table[keyIndex].headPair = h.table[keyIndex].headPair.nextPair
		}
		return nil
	} else {
		updateCount := 0
		for h.table[keyIndex].headPair != nil {
			if h.table[keyIndex].headPair.key == key {
				h.table[keyIndex].headPair.value = value
				updateCount++
			}
			h.table[keyIndex].headPair = h.table[keyIndex].headPair.nextPair
		}
		if updateCount > 0 {
			return nil
		} else {
			return errors.New("key doesn't already exist")
		}
	}
}

// Delete deletes a key/value pair from the hashtable.
// Should return an error if the given key doesn't exist.
func (h *Hashtable) Delete(key string) error {
	// TODO: Implement
	keyIndex := generateKeyIndex(key)
	if h.table[keyIndex].key == key {
		h.table[keyIndex].key = ""
		h.table[keyIndex].value = 0
		for h.table[keyIndex].headPair != nil {
			if h.table[keyIndex].headPair.key == key {
				h.table[keyIndex].headPair.key = ""
				h.table[keyIndex].headPair.value = 0
			}
			h.table[keyIndex].headPair = h.table[keyIndex].headPair.nextPair
		}
		return nil
	} else {
		deleteCount := 0
		for h.table[keyIndex].headPair != nil {
			if h.table[keyIndex].headPair.key == key {
				h.table[keyIndex].headPair.key = ""
				h.table[keyIndex].headPair.value = 0
				deleteCount++
			}
			h.table[keyIndex].headPair = h.table[keyIndex].headPair.nextPair
		}
		if deleteCount > 0 {
			return nil
		} else {
			return errors.New("key doesn't already exist")
		}
	}
}

// Exists returns true if the key exists in the hashtable, false otherwise.
func (h *Hashtable) Exists(key string) bool {
	// TODO: Implement
	keyIndex := generateKeyIndex(key)
	if h.table[keyIndex].key == key {
		return true
	} else {
		for h.table[keyIndex].headPair != nil {
			if h.table[keyIndex].headPair.key == key {
				return true
			}
			h.table[keyIndex].headPair = h.table[keyIndex].headPair.nextPair
		}
		return false
	}
}

// Get returns the value associated with the given key.
// Should return an error if value doesn't exist.
func (h *Hashtable) Get(key string) (int, error) {
	// TODO: Implement
	keyIndex := generateKeyIndex(key)
	if h.table[keyIndex].key == key {
		return h.table[keyIndex].value, nil
	} else {
		for h.table[keyIndex].headPair != nil {
			if h.table[keyIndex].headPair.key == key {
				return h.table[keyIndex].headPair.value, nil
			}
			h.table[keyIndex].headPair = h.table[keyIndex].headPair.nextPair
		}
		return 0, errors.New("key doesn't exist")
	}
}

func generateKeyIndex(key string) uint16 {
	hs := sha256.Sum256([]byte(key))
	return (uint16(hs[0]) << 8) + uint16(hs[1])
}
