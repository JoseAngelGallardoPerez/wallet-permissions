<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddClientDefaultGroup extends Migration
{
    const CLIENTS_DEFAULT_GROUP = "Clients Default Group";

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        DB::table("permissions_groups")->insert([
            "id" => 2,
            "name" => self::CLIENTS_DEFAULT_GROUP,
            "description" => "Default group for users with role `client`",
            "scope" => "client",
            "created_at" => date("Y-m-d H:i:s"),
            "updated_at" => date("Y-m-d H:i:s")
        ]);

        $actions = [
            "view_user_profiles",
            "view_modify_settings"
        ];

        $group = $this->selectGroupByName(self::CLIENTS_DEFAULT_GROUP);
        foreach ($actions as $key) {
            $action = $this->selectPermissionActionByKey($key);

            DB::insert('INSERT INTO permissions_groups_actions (`group_id`, `action_id`) VALUES (?, ?)',
                [$group->id, $action->id]);
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        DB::delete("DELETE FROM permissions_groups WHERE `name` = ?", [self::CLIENTS_DEFAULT_GROUP]);
    }

    private function selectPermissionActionByKey(string $key)
    {
        $actions = DB::select("SELECT * FROM permissions_actions WHERE `key` = ?", [$key]);
        if (empty($actions)) {
            return [];
        }
        return $actions[0];
    }

    private function selectGroupByName(string $groupName)
    {
        $group = DB::select("SELECT * FROM permissions_groups WHERE `name` = ?", [$groupName]);
        if (empty($group)) {
            throw new \Exception("group with name $groupName not found");
        }
        return $group[0];
    }
}
