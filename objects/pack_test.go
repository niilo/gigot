// Copyright 2012 Rémy Oudompheng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objects

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFindInPack(t *testing.T) {
	// Take a random packfile in our own repository.
	packs, err := filepath.Glob("../.git/objects/pack/pack-*.pack")
	if err != nil || len(packs) == 0 {
		t.Fatalf("globbing failed: %v", err, packs)
	}
	pname := packs[0]
	pack, err := os.Open(pname)
	if err != nil {
		t.Fatalf("count not open %s: %s", pname, err)
	}
	idx, err := os.Open(pname[:len(pname)-5] + ".idx")
	if err != nil {
		t.Fatal(err)
	}
	packstat, _ := pack.Stat()
	t.Logf("opening pack %s (%d bytes)", pname, packstat.Size())
	pk, err := NewPackReader(pack, idx)
	if err != nil {
		t.Fatal(err)
	}

	// Take the object ID of a random ref.
	refs, err := ioutil.ReadFile("../.git/info/refs")
	if err != nil {
		t.Fatal(err)
	}
	id := bytes.Fields(refs)[0]
	if len(id) != 40 {
		t.Fatal("invalid commit ID %q in info/refs", id)
	}
	var refhash [20]byte
	hex.Decode(refhash[:], id)
	t.Logf("lookup %040x", refhash)
	off, err := pk.findObject(refhash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("offset=%v", off)
}
