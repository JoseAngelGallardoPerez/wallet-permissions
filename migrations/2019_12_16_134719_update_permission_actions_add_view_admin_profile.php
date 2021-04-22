<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Support\Facades\DB;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class UpdatePermissionActionsAddViewAdminProfile extends Migration
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
        $this->removeOldPermissions();

        $this->updateParentsAndCategory();
        $this->updatePermissionsGroupActions();
    }

    private function removeOldPermissions()
    {
        DB::delete('DELETE FROM permissions_actions WHERE `key` = ?', ['create_view_admin_profiles']);
    }

    private function updateParentsAndCategory()
    {
        $permissions = [
            // Profiles
            'create_profile' => [
                'category' => null,
                'parent' => 'view_user_profiles',
                'sort' => 1000,
                'isHidden' => false,
                'newName' => 'Create User Profiles'
            ],

            'view_admin_profiles' => [
                'category' => 'Profiles',
                'parent' => null,
                'sort' => 800,
                'isHidden' => false,
                'name' => 'View Admin Profiles'
            ],
            'create_admin_profiles' => [
                'category' => null,
                'parent' => 'view_admin_profiles',
                'sort' => 900,
                'isHidden' => false,
                'name' => 'Create Admin Profiles'
            ],
            'modify_admin_profiles' => [
                'category' => null,
                'parent' => 'view_admin_profiles',
                'sort' => 850,
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

            $action = $this->selectPermissionActionByKey($key);
            if ($action) {
                DB::update('UPDATE permissions_actions SET `parent_id` = ?, `category_id`= ?, `is_hidden`=?, `sort`=?  WHERE `key` = ?',
                    [$parentId, $categoryId, (int)$permission['isHidden'], $permission['sort'], $key]);
                if (isset($permission['newName'])) {
                    DB::update('UPDATE permissions_actions SET `name` = ? WHERE `key` = ?', [$permission['newName'], $key]);
                }
            } else {
                $this->createAction($key, $permission, $parentId);
            }
        }
    }

    private function createAction(string $key, array $action, $parentId)
    {
        $categoryId = null;
        if (!empty($action['category'])) {
            $category = $this->selectCategoryByName($action['category']);
            $categoryId = $category->id;
        }
        DB::insert('INSERT INTO permissions_actions (`key`, `name`, `parent_id`, `category_id`, `sort`) VALUES (?, ?, ?, ?, ?)',
            [$key, $action['name'], $parentId, $categoryId, $action['sort']]);
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
            return [];
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
