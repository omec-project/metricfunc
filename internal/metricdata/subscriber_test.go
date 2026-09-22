// Copyright (c) 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"testing"

	"github.com/omec-project/util/metricinfo"
)

// resetSubscriberState clears the package-level subscriber stores so each test
// starts from a clean slate, since metricData, smSessCounts and subPresences are
// package-level state shared across the test binary.
func resetSubscriberState() {
	metricData.Subscribers = make(map[string]*metricinfo.CoreSubscriber)
	smSessCounts = map[smSessKey]uint64{}
	subPresences = map[string]*smSessPresence{}
}

func TestUpdateSubscriber_UncountedMigration_DoesNotCorruptSharedTuple(t *testing.T) {
	resetSubscriberState()

	sharedTuple := metricinfo.CoreSubscriber{SmfIp: "10.0.0.1", Slice: "slice1", Dnn: "internet", UpfName: "upf1"}

	// A live SMF session already owns the shared tuple.
	live := sharedTuple
	live.Imsi = "imsi-live"
	HandleSubscriberEvent(&metricinfo.CoreSubscriberData{Subscriber: live, Operation: metricinfo.SubsOpAdd}, metricinfo.NfTypeSmf)

	if got := smSessCounts[smSessKeyOf(&live)]; got != 1 {
		t.Fatalf("expected live session tuple count 1, got %d", got)
	}

	// An AMF-added subscriber whose record happens to carry the same tuple values,
	// but was never counted because it was AMF-sourced.
	uncounted := sharedTuple
	uncounted.Imsi = "imsi-uncounted"
	HandleSubscriberEvent(&metricinfo.CoreSubscriberData{Subscriber: uncounted, Operation: metricinfo.SubsOpAdd}, metricinfo.NfTypeAmf)

	if p := subPresences[uncounted.Imsi]; p == nil || p.smf {
		t.Fatalf("expected AMF-added subscriber to be uncounted for smf, got %+v", p)
	}

	// An SMF mod migrates it to a new tuple. Since the old tuple was never counted for
	// this subscriber, this must not decrement the still-live shared tuple's count.
	migrated := metricinfo.CoreSubscriber{Imsi: uncounted.Imsi, SmfIp: "10.0.0.2", Slice: "slice2", Dnn: "ims", UpfName: "upf2"}
	HandleSubscriberEvent(&metricinfo.CoreSubscriberData{Subscriber: migrated, Operation: metricinfo.SubsOpMod}, metricinfo.NfTypeSmf)

	if got := smSessCounts[smSessKeyOf(&live)]; got != 1 {
		t.Errorf("shared tuple count corrupted by unrelated migration: got %d, want 1", got)
	}
	if got := smSessCounts[smSessKeyOf(&migrated)]; got != 1 {
		t.Errorf("expected migrated tuple count 1, got %d", got)
	}
}

func TestDeleteSubscriber_AmfDelete_DoesNotDecrementLiveSmfSession(t *testing.T) {
	resetSubscriberState()

	sub := metricinfo.CoreSubscriber{Imsi: "imsi-1", SmfIp: "10.0.0.1", Slice: "slice1", Dnn: "internet", UpfName: "upf1"}
	HandleSubscriberEvent(&metricinfo.CoreSubscriberData{Subscriber: sub, Operation: metricinfo.SubsOpAdd}, metricinfo.NfTypeSmf)

	key := smSessKeyOf(&sub)
	if got := smSessCounts[key]; got != 1 {
		t.Fatalf("expected session count 1 after add, got %d", got)
	}

	// AMF tears down the UE context first, without the SMF PDU session ending.
	HandleSubscriberEvent(&metricinfo.CoreSubscriberData{Subscriber: sub, Operation: metricinfo.SubsOpDel}, metricinfo.NfTypeAmf)

	if got := smSessCounts[key]; got != 1 {
		t.Errorf("AMF delete must not touch the live SMF session count: got %d, want 1", got)
	}
	if _, ok := metricData.Subscribers[sub.Imsi]; !ok {
		t.Errorf("subscriber record must survive an AMF-only delete while SMF session is live")
	}

	// SMF later confirms the PDU session actually ended.
	HandleSubscriberEvent(&metricinfo.CoreSubscriberData{Subscriber: sub, Operation: metricinfo.SubsOpDel}, metricinfo.NfTypeSmf)

	if _, ok := smSessCounts[key]; ok {
		t.Errorf("expected session series to be removed after the SMF-sourced delete")
	}
	if _, ok := metricData.Subscribers[sub.Imsi]; ok {
		t.Errorf("subscriber record should be removed once both sources report delete")
	}
}
