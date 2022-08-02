package ee14

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/flavours/postgres14/common14"
	"github.com/wader/fq/pkg/decode"
)

// type = struct PageHeaderData {
/*    0      |     8 */ // PageXLogRecPtr pd_lsn;
/*    8      |     2 */ // uint16 pd_checksum;
/*   10      |     2 */ // uint16 pd_flags;
/*   12      |     2 */ // LocationIndex pd_lower;
/*   14      |     2 */ // LocationIndex pd_upper;
/*   16      |     2 */ // LocationIndex pd_special;
/*   18      |     2 */ // uint16 pd_pagesize_version;
/*   20      |     0 */ // ItemIdData pd_linp[];
func DecodePageHeaderData(d *decode.D) {
	heap := common14.GetHeapD(d)
	page := heap.Page

	d.FieldStruct("pd_lsn", func(d *decode.D) {
		/*    0      |     4 */ // uint32 xlogid;
		/*    4      |     4 */ // uint32 xrecoff;
		d.FieldU32("xlogid", common.HexMapper)
		d.FieldU32("xrecoff", common.HexMapper)
	})
	d.FieldU16("pd_checksum")
	d.FieldU16("pd_flags")
	page.PdLower = uint16(d.FieldU16("pd_lower"))
	page.PdUpper = uint16(d.FieldU16("pd_upper"))
	page.PdSpecial = uint16(d.FieldU16("pd_special"))
	page.PdPagesizeVersion = uint16(d.FieldU16("pd_pagesize_version"))

	// ItemIdData pd_linp[];
	page.ItemsEnd = int64(page.PagePosBegin*8) + int64(page.PdLower*8)
	d.FieldArray("pd_linp", common14.DecodeItemIds)
}

// type = struct HeapPageSpecialData {
/*    0      |     8 */ // TransactionId pd_xid_base;
/*    8      |     8 */ // TransactionId pd_multi_base;
/*   16      |     4 */ // ShortTransactionId pd_prune_xid;
/*   20      |     4 */ // uint32 pd_magic;
//
/* total size (bytes):   24 */
func DecodePageSpecial(d *decode.D) {
	heap := common14.GetHeapD(d)
	page := heap.Page

	specialPos := int64(page.PagePosBegin*8) + int64(page.PdSpecial*8)
	d.SeekAbs(specialPos)

	d.FieldStruct("HeapPageSpecialData", func(d *decode.D) {
		page.PdXidBase = d.FieldU64("pd_xid_base")
		page.PdMultiBase = d.FieldU64("pd_multi_base")
		page.PdPruneXid = d.FieldU32("pd_prune_xid")
		page.PdMagic = d.FieldU32("pd_magic")
	})
}
