<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class UpdatePermissionActions extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $this->updateParentsAndCategory();
        $this->updatePermissionsGroupActions();
    }

    private function updateParentsAndCategory()
    {
        $permissions = [
            // Requests
            'execute_cancel_pending_transfer_requests' => [
                'category' => 'Requests',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],
            'approve_cancel_registration_requests' => [
                'category' => 'Requests',
                'parent' => null,
                'sort' => 900,
                'isHidden' => false,
            ],
            'import_transfer_request_updates' => [
                'category' => 'Requests',
                'parent' => null,
                'sort' => 800,
                'isHidden' => false,
            ],

            // System Log
            'view_system_log' => [
                'category' => 'System Log',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],

            // Accounts
            'view_accounts' => [
                'category' => 'Accounts',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],
            'create_accounts' => [
                'category' => null,
                'parent' => 'view_accounts',
                'sort' => 900,
                'isHidden' => false,
            ],
            'create_accounts_with_initial_balance' => [
                'category' => null,
                'parent' => 'create_accounts',
                'sort' => 800,
                'isHidden' => false,
            ],
            'modify_accounts' => [
                'category' => null,
                'parent' => 'view_accounts',
                'sort' => 850,
                'isHidden' => false,
            ],
            'manual_debit_credit_accounts' => [
                'category' => null,
                'parent' => 'modify_accounts',
                'sort' => 750,
                'isHidden' => false,
            ],
            'modify_cards' => [
                'category' => 'Accounts',
                'parent' => null,
                'sort' => 850,
                'isHidden' => false,
            ],

            'view_revenue' => [
                'category' => 'Accounts',
                'parent' => null,
                'sort' => 800,
                'isHidden' => false,
            ],
            'manage_revenue' => [
                'category' => null,
                'parent' => 'view_revenue',
                'sort' => 1000,
                'isHidden' => false,
            ],
            'create_view_cards' => [ // hidden in the front
                'category' => 'Accounts',
                'parent' => null,
                'sort' => 100,
                'isHidden' => true,
            ],

            // Profiles
            'create_profile' => [
                'category' => null,
                'parent' => 'view_user_profiles',
                'sort' => 1000,
                'isHidden' => false,
            ],
            'view_user_profiles' => [
                'category' => 'Profiles',
                'parent' => null,
                'sort' => 900,
                'isHidden' => false,
            ],
            'modify_user_profiles' => [
                'category' => null,
                'parent' => 'view_user_profiles',
                'sort' => 900,
                'isHidden' => false,
            ],
            'modify_admin_profiles' => [
                'category' => null,
                'parent' => 'view_user_profiles',
                'sort' => 850,
                'isHidden' => false,
            ],
            'create_view_admin_profiles' => [ // hidden in the front
                'category' => null,
                'parent' => 'view_user_profiles',
                'sort' => 100,
                'isHidden' => true,
            ],

            // Transfers
            'initiate_execute_user_transfers' => [
                'category' => 'Transfers',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],

            // Messages
            'send_reply_internal_messages' => [
                'category' => 'Messages',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],

            // Settings
            'view_modify_settings' => [
                'category' => 'Settings',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],
            'create_modify_iwt_bank_accounts' => [
                'category' => null,
                'parent' => 'view_modify_settings',
                'sort' => 900,
                'isHidden' => false,
            ],
            'modify_account_types' => [
                'category' => null,
                'parent' => 'view_modify_settings',
                'sort' => 800,
                'isHidden' => false,
                'newName' => 'Create/Modify Account types',
            ],

            // Reports
            'view_user_reports' => [
                'category' => 'Reports',
                'parent' => null,
                'sort' => 1000,
                'isHidden' => false,
            ],
            'view_general_system_reports' => [
                'category' => 'Reports',
                'parent' => null,
                'sort' => 900,
                'isHidden' => false,
            ],
        ];

        foreach ($permissions as $key => $permission) {
            $parentId = null;
            $categoryId = null;
            if (!empty($permission['parent'])) {
                $parentAction = $this->selectPermissionActionByKey($permission['parent']);
                $parentId = $parentAction->id;
            }
            if (!empty($permission['category'])) {
                $category = $this->selectCategoryByName($permission['category']);
                $categoryId = $category->id;
            }

            DB::update('UPDATE permissions_actions SET `parent_id` = ?, `category_id`= ?, `is_hidden`=?, `sort`=?  WHERE `key` = ?', [$parentId, $categoryId, (int)$permission['isHidden'], $permission['sort'], $key]);
            if (isset($permission['newName'])) {
                DB::update('UPDATE permissions_actions SET `name` = ? WHERE `key` = ?', [$permission['newName'], $key]);
            }
        }
    }

    /**
     * if a group has child permission then it must have parent permission
     */
    private function updatePermissionsGroupActions()
    {
        $permissionsGroups = DB::select('SELECT * FROM permissions_groups_actions');
        foreach ($permissionsGroups as $permGroup) {
            $action = $this->selectPermissionActionById($permGroup->action_id);
            $this->addPermissionGroupActionIfNeed($action, $permGroup->group_id);
        }
    }

    /**
     * add to permissions_groups_actions new records if permission_actions has a parent.
     *  the method checks all parents recursive
     *
     * @param stdClass $action
     * @param int $groupId
     * @throws Exception
     */
    private function addPermissionGroupActionIfNeed(\stdClass $action, int $groupId)
    {
        if ($action->parent_id) {
            // check if exists parent permission for group
            $sql = 'SELECT * FROM permissions_groups_actions WHERE group_id = ? AND action_id = ?';
            $permissionsGroupsParent = DB::select($sql, [$groupId, $action->parent_id]);
            if (empty($permissionsGroupsParent)) {
                DB::insert('INSERT INTO permissions_groups_actions (`group_id`, `action_id`) VALUES (?, ?)',
                    [$groupId, $action->parent_id]);
            }
            $parentAction = $this->selectPermissionActionById($action->parent_id);
            $this->addPermissionGroupActionIfNeed($parentAction, $groupId);
        }
    }

    private function selectPermissionActionById(int $id)
    {
        $actions = DB::select('SELECT * FROM permissions_actions WHERE `id` = ?', [$id]);
        if (empty($actions)) {
            throw new \Exception("permissions_actions with id $id not found");
        }
        return $actions[0];
    }

    private function selectPermissionActionByKey(string $key)
    {
        $actions = DB::select('SELECT * FROM permissions_actions WHERE `key` = ?', [$key]);
        if (empty($actions)) {
            throw new \Exception("permissions_actions with key $key not found");
        }

        return $actions[0];
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
