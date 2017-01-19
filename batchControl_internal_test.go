// Copyright 2017 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"strings"
	"testing"
)

// TestParseBatchControl parses a known Batch ControlRecord string.
func TestParseBatchControl(t *testing.T) {
	var line = "82250000010005320001000000010500000000000000origid                             076401250000001"
	r := NewReader(strings.NewReader(line))
	r.line = line
	r.currentBatch.Header.BatchNumber = 1
	r.currentBatch.Header.ServiceClassCode = 225
	r.currentBatch.Header.CompanyIdentification = "origid"
	r.currentBatch.addEntryDetail(EntryDetail{TransactionCode: 27, Amount: 10500, RDFIIdentification: 5320001})
	err := r.parseBatchControl()
	if err != nil {
		t.Errorf("unknown error: %v", err)
	}
	record := r.file.Batches[0].Control

	if record.recordType != "8" {
		t.Errorf("RecordType Expected '8' got: %v", record.recordType)
	}
	if record.ServiceClassCode != 225 {
		t.Errorf("ServiceClassCode Expected '225' got: %v", record.ServiceClassCode)
	}
	if record.EntryAddendaCountField() != "000001" {
		t.Errorf("EntryAddendaCount Expected '000001' got: %v", record.EntryAddendaCountField())
	}
	if record.EntryHashField() != "0005320001" {
		t.Errorf("EntryHash Expected '0005320001' got: %v", record.EntryHashField())
	}
	if record.TotalDebitEntryDollarAmountField() != "000000010500" {
		t.Errorf("TotalDebitEntryDollarAmount Expected '000000010500' got: %v", record.TotalDebitEntryDollarAmountField())
	}
	if record.TotalCreditEntryDollarAmountField() != "000000000000" {
		t.Errorf("TotalCreditEntryDollarAmount Expected '000000000000' got: %v", record.TotalCreditEntryDollarAmountField())
	}
	if record.CompanyIdentificationField() != "origid    " {
		t.Errorf("CompanyIdentification Expected 'origid    ' got: %v", record.CompanyIdentificationField())
	}
	if record.MessageAuthenticationCodeField() != "                   " {
		t.Errorf("MessageAuthenticationCode Expected '                   ' got: %v", record.MessageAuthenticationCodeField())
	}
	if record.reserved != "      " {
		t.Errorf("Reserved Expected '      ' got: %v", record.reserved)
	}
	if record.ODFIIdentificationField() != "07640125" {
		t.Errorf("OdfiIdentification Expected '07640125' got: %v", record.ODFIIdentificationField())
	}
	if record.BatchNumberField() != "0000001" {
		t.Errorf("BatchNumber Expected '0000001' got: %v", record.BatchNumberField())
	}
}

// TestBCString validats that a known parsed file can be return to a string of the same value
func TestBCString(t *testing.T) {
	var line = "82250000010005320001000000010500000000000000origid                             076401250000001"
	r := NewReader(strings.NewReader(line))
	r.line = line
	r.currentBatch.Header.BatchNumber = 1
	r.currentBatch.Header.ServiceClassCode = 225
	r.currentBatch.Header.CompanyIdentification = "origid"
	r.currentBatch.addEntryDetail(EntryDetail{TransactionCode: 27, Amount: 10500, RDFIIdentification: 5320001})
	err := r.parseBatchControl()
	if err != nil {
		t.Errorf("unknown error: %v", err)
	}
	record := r.file.Batches[0].Control

	if record.String() != line {
		t.Errorf("Strings do not match")
	}
}

// TestValidateBCRecordType ensure error if recordType is not 8
func TestValidateBCRecordType(t *testing.T) {
	bc := NewBatchControl()
	bc.recordType = "2"
	if err := bc.Validate(); err != nil {
		if !strings.Contains(err.Error(), ErrRecordType.Error()) {
			t.Errorf("Expected RecordType Error got: %v", err)
		}
	}
}
