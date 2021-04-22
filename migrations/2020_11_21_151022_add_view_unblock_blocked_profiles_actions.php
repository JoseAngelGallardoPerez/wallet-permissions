<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddViewUnblockBlockedProfilesActions extends Migration
{
    const TABLE = "permissions_actions";
    const CATEGORY_PROFILES = "Profiles";

    public function up()
    {
        $category = $this->selectCategoryByName(self::CATEGORY_PROFILES);
        $categoryId = $category->id;

        DB::table(self::TABLE)->insert([
            "name" => "View, Unblock Blocked Profiles",
            "key" => "view_unblock_blocked_profiles",
            "description" => null,
            "sort" => 860,
            "is_hidden" => false,
            "parent_id" => null,
            "category_id" => $categoryId
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
            "view_unblock_blocked_profiles"
        ];
        foreach ($keys as $key) {
            DB::delete("DELETE FROM permissions_actions WHERE `key` = ?", [$key]);
        }
    }

    private function selectCategoryByName(string $categoryName)
    {
        $categories = DB::select('SELECT * FROM permissions_categories WHERE `name` = ?', [$categoryName]);
        if (empty($categories)) {
            throw new \Exception("category with name $categoryName not found");
        }
        return $categories[0];
    }
}
