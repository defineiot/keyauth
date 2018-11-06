package mysql_test

import (
	"testing"
)

func TestRole(t *testing.T) {
	t.Run("CreateRoleOK", testCreateRoleOK)
	t.Run("ListRoleOK", testListRoleOK)
	t.Run("GetRoleOK", testGetRoleOK)
	t.Run("DeleteRoleOK", testDeleteRoleOK)
	t.Run("LinkFeatruesOK", testLinkFeatures)
	t.Run("UnlinkFeaturesOK", testUnlinkFeatures)
}

func testCreateRoleOK(t *testing.T) {

}

func testListRoleOK(t *testing.T) {

}

func testGetRoleOK(t *testing.T) {

}

func testDeleteRoleOK(t *testing.T) {

}

func testLinkFeatures(t *testing.T) {

}

func testUnlinkFeatures(t *testing.T) {

}
