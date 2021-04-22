<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSettingsActions extends Migration
{
    const TABLE = "permissions_actions";
    const CATEGORY_SETTINGS = "Settings";

    public function up()
    {
        $this->renameSettingPermissionKey('view_modify_settings', 'view_settings', 'View Settings');

        $parent = $this->selectPermissionActionByKey('view_settings');
        $parentId = $parent->id;

        $category = $this->selectCategoryByName(self::CATEGORY_SETTINGS);
        $categoryId = $category->id;

        DB::table(self::TABLE)->insert([
            "name" => "Create Settings",
            "key" => "create_settings",
            "description" => null,
            "sort" => 1010,
            "is_hidden" => false,
            "parent_id" => $parentId,
            "category_id" => null
        ]);

        DB::table(self::TABLE)->insert([
            "name" => "Modify Settings",
            "key" => "modify_settings",
            "description" => null,
            "sort" => 1011,
            "is_hidden" => false,
            "parent_id" => $parentId,
            "category_id" => null
        ]);

        DB::table(self::TABLE)->insert([
            "name" => "Remove Settings",
            "key" => "remove_settings",
            "description" => null,
            "sort" => 1012,
            "is_hidden" => false,
            "parent_id" => $parentId,
            "category_id" => null
        ]);
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        $keys = [
            "remove_settings",
            "modify_settings",
            "create_settings"
        ];
        foreach ($keys as $key) {
            DB::delete("DELETE FROM permissions_actions WHERE `key` = ?", [$key]);
        }

        $this->renameSettingPermissionKey('view_settings', 'view_modify_settings', 'View/Modify Settings');
    }

    private function selectCategoryByName(string $categoryName)
    {
        $categories = DB::select('SELECT * FROM permissions_categories WHERE `name` = ?', [$categoryName]);
        if (empty($categories)) {
            throw new \Exception("category with name $categoryName not found");
        }
        return $categories[0];
    }

    private function renameSettingPermissionKey($oldKey, $newKey, $newName)
    {
        DB::delete('UPDATE permissions_actions SET `key`=?, `name`=? WHERE `key` = ?', [$newKey, $newName, $oldKey]);
    }

    private function selectPermissionActionByKey(string $key)
    {
        $actions = DB::select("SELECT * FROM permissions_actions WHERE `key` = ?", [$key]);
        if (empty($actions)) {
            return [];
        }
        return $actions[0];
    }
}
