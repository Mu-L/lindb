package metadb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

var testPath = "test"

func TestMetadataBackend_new(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkDir = fileutil.MkDirIfNotExist
		nsBucketName = []byte("ns")
		metricBucketName = []byte("m")
	}()

	// test: new success
	db, err := newMetadataBackend(testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// test: can't re-open
	db1, err := newMetadataBackend(testPath)
	assert.Error(t, err)
	assert.Nil(t, db1)

	// close db
	err = db.Close()
	assert.NoError(t, err)

	// test: create namespace bucket err
	nsBucketName = []byte("")
	db1, err = newMetadataBackend(testPath)
	assert.Error(t, err)
	assert.Nil(t, db1)

	// test: create metric bucket err
	nsBucketName = []byte("ns")
	metricBucketName = []byte("")
	db1, err = newMetadataBackend(testPath)
	assert.Error(t, err)
	assert.Nil(t, db1)

	// test: create parent path err
	mkDir = func(path string) error {
		return fmt.Errorf("err")
	}
	db, err = newMetadataBackend(testPath)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestMetadataBackend_suggestNamespace(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)

	values, err := db.suggestNamespace("ns", 100)
	assert.Equal(t, []string{"ns-1", "ns-2"}, values)
	assert.NoError(t, err)

	values, err = db.suggestNamespace("ns-2", 100)
	assert.Equal(t, []string{"ns-2"}, values)
	assert.NoError(t, err)

	values, err = db.suggestNamespace("ns", 1)
	assert.Equal(t, []string{"ns-1"}, values)
	assert.NoError(t, err)

	values, err = db.suggestNamespace("aans", 1)
	assert.Empty(t, values)
	assert.NoError(t, err)
}

func TestMetadataBackend_suggestMetricName(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)

	values, err := db.suggestMetricName("ns-3", "name", 100)
	assert.Empty(t, values)
	assert.NoError(t, err)

	values, err = db.suggestMetricName("ns-2", "name", 100)
	assert.Equal(t, []string{"name2", "name3"}, values)
	assert.NoError(t, err)

	values, err = db.suggestMetricName("ns-2", "name", 1)
	assert.Equal(t, []string{"name2"}, values)
	assert.NoError(t, err)

	values, err = db.suggestMetricName("ns-2", "name3", 1)
	assert.Equal(t, []string{"name3"}, values)
	assert.NoError(t, err)
}

func TestMetadataBackend_gen_id(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := newMockMetadataBackend(t)
	assert.Equal(t, uint32(1), db.genMetricID())
	assert.Equal(t, uint32(2), db.genMetricID())
	assert.Equal(t, uint32(1), db.genTagKeyID())
	assert.Equal(t, uint32(2), db.genTagKeyID())

	event := mockMetadataEvent()
	// save metadata
	err := db.saveMetadata(event)
	assert.NoError(t, err)
	err = db.Close()
	assert.NoError(t, err)
	// re-open,load new tag key/metric id sequence
	db = newMockMetadataBackend(t)
	assert.Equal(t, uint32(5), db.genMetricID())
	assert.Equal(t, uint32(5), db.genTagKeyID())
}

func TestMetadataBackend_loadMetricMetadata(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)
	_, err := db.loadMetricMetadata("ns1", "name2")
	assert.Equal(t, series.ErrNotFound, err)

	meta, err := db.loadMetricMetadata("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), meta.getMetricID())
	assert.Equal(t, []tag.Meta{{Key: "tagKey-2", ID: 4}, {Key: "tagKey-3", ID: 3}}, meta.getAllTagKeys())
	assert.Equal(t, []field.Meta{
		{ID: 1, Name: "f3", Type: field.MaxField},
		{ID: 3, Name: "f4", Type: field.SumField},
	}, meta.getAllFields())
	m := meta.(*metricMetadata)
	assert.Equal(t, int32(3), m.fieldIDSeq.Load())
	fID, err := meta.createField("f5", field.SumField)
	assert.Equal(t, uint16(4), fID)
	assert.NoError(t, err)

	// test: metric id not exist
	_, err = db.getMetricMetadata(999)
	assert.Error(t, err)
}

func TestMetadataBackend_getTagKeyID(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)
	metricID, _ := db.getMetricID("ns-1", "name2")
	_, err := db.getTagKeyID(metricID, "ggg")
	assert.Equal(t, series.ErrNotFound, err)
	_, err = db.getTagKeyID(99, "tagKey-3")
	assert.Equal(t, series.ErrNotFound, err)

	tagKeyID, err := db.getTagKeyID(metricID, "tagKey-3")
	assert.NoError(t, err)
	assert.Equal(t, uint32(3), tagKeyID)
}

func TestMetadataBackend_getAllTagKeys(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)
	_, err := db.getAllTagKeys(88)
	assert.Equal(t, series.ErrNotFound, err)
	values, err := db.getAllTagKeys(2)
	assert.NoError(t, err)
	assert.Equal(t, []tag.Meta{{Key: "tagKey-2", ID: 4}, {Key: "tagKey-3", ID: 3}}, values)
}

func TestMetadataBackend_getField(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)
	_, err := db.getField(99, "f3")
	assert.Equal(t, series.ErrNotFound, err)
	_, err = db.getField(2, "f33")
	assert.Equal(t, series.ErrNotFound, err)
	f, err := db.getField(2, "f3")
	assert.NoError(t, err)
	assert.Equal(t, field.Meta{ID: 1, Name: "f3", Type: field.MaxField}, f)
}

func TestMetadataBackend_getAllFields(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := mockMetadataBackend(t)
	_, err := db.getAllFields(99)
	assert.Equal(t, series.ErrNotFound, err)
	fields, err := db.getAllFields(2)
	assert.Equal(t, []field.Meta{
		{ID: 1, Name: "f3", Type: field.MaxField},
		{ID: 3, Name: "f4", Type: field.SumField},
	}, fields)
	assert.NoError(t, err)
}

func TestMetadataBackend_saveMetadata(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := newMockMetadataBackend(t)
	event := mockMetadataEvent()
	err := db.saveMetadata(event)
	assert.NoError(t, err)
	// save duplicate event
	err = db.saveMetadata(event)
	assert.NoError(t, err)

	metricID, err := db.getMetricID("ns-1", "name1")
	assert.Equal(t, uint32(1), metricID)
	assert.NoError(t, err)
	metricID, err = db.getMetricID("ns-2", "name3")
	assert.Equal(t, uint32(3), metricID)
	assert.NoError(t, err)

	_, err = db.getMetricID("ns-2", "name5")
	assert.Equal(t, series.ErrNotFound, err)
}

func TestMetadataBackend_save_err(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		tagBucketName = []byte("t")
		fieldBucketName = []byte("f")
	}()
	db := newMockMetadataBackend(t)
	// ns is empty
	e := newMetadataUpdateEvent()
	e.addMetric("", "name1", 1)
	err := db.saveMetadata(e)
	assert.Error(t, err)

	// metric name is empty
	e = newMetadataUpdateEvent()
	e.addMetric("ns-2", "", 1)
	err = db.saveMetadata(e)
	assert.Error(t, err)

	// tag key is empty
	e = newMetadataUpdateEvent()
	e.addTagKey(1, tag.Meta{Key: "", ID: 1})
	err = db.saveMetadata(e)
	assert.Error(t, err)

	// field name is empty
	e = newMetadataUpdateEvent()
	e.addField(1, field.Meta{ID: 1, Name: "", Type: field.SummaryField})
	err = db.saveMetadata(e)
	assert.Error(t, err)

	// tag key bucket name is empty
	tagBucketName = []byte("")
	e = newMetadataUpdateEvent()
	e.addTagKey(1, tag.Meta{Key: "empty_tag_key", ID: 1})
	err = db.saveMetadata(e)
	assert.Error(t, err)

	// field bucket name is empty
	tagBucketName = []byte("t")
	fieldBucketName = []byte("")
	e = newMetadataUpdateEvent()
	e.addField(1, field.Meta{ID: 1, Name: "", Type: field.SummaryField})
	err = db.saveMetadata(e)
	assert.Error(t, err)
}

func TestMetadataBackend_sync(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := newMockMetadataBackend(t)
	err := db.sync()
	assert.NoError(t, err)
	err = db.Close()
	assert.NoError(t, err)
}

func newMockMetadataBackend(t *testing.T) MetadataBackend {
	db, err := newMetadataBackend(testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	return db
}

func mockMetadataBackend(t *testing.T) MetadataBackend {
	db := newMockMetadataBackend(t)
	event := mockMetadataEvent()
	err := db.saveMetadata(event)
	assert.NoError(t, err)
	return db
}

func mockMetadataEvent() *metadataUpdateEvent {
	e := newMetadataUpdateEvent()
	e.addMetric("ns-1", "name1", 1)
	e.addMetric("ns-1", "name2", 2)
	e.addMetric("ns-2", "name3", 3)
	e.addMetric("ns-2", "name2", 4)

	// tags
	e.addTagKey(1, tag.Meta{Key: "tagKey-1", ID: 1})
	e.addTagKey(1, tag.Meta{Key: "tagKey-2", ID: 2})
	e.addTagKey(2, tag.Meta{Key: "tagKey-3", ID: 3})
	e.addTagKey(2, tag.Meta{Key: "tagKey-2", ID: 4})

	// fields
	e.addField(1, field.Meta{ID: 1, Name: "f1", Type: field.SummaryField})
	e.addField(1, field.Meta{ID: 2, Name: "f2", Type: field.MinField})
	e.addField(2, field.Meta{ID: 1, Name: "f3", Type: field.MaxField})
	e.addField(2, field.Meta{ID: 3, Name: "f4", Type: field.SumField})

	return e
}