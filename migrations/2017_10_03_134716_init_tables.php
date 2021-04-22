<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class InitTables extends Migration
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
        // skip the migration if there are another migrations
        // It means this migration was already applied
        $migrations = DB::select('SELECT * FROM migrations LIMIT 1');
        if (!empty($migrations)) {
            return;
        }
        $oldMigrationTable = DB::select("SHOW TABLES LIKE 'schema_migrations'");
        if (!empty($oldMigrationTable)) {
            return;
        }

        DB::beginTransaction();

        try {
            app("db")->getPdo()->exec($this->getSql());
        } catch (\Throwable $e) {
            DB::rollBack();
            throw $e;
        }

        DB::commit();
    }

    private function getSql()
    {
        return <<<SQL
            CREATE TABLE `permissions_actions` (
              `id` int(10) UNSIGNED NOT NULL,
              `name` varchar(128) NOT NULL,
              `key` varchar(128) NOT NULL,
              `description` varchar(255) DEFAULT NULL,
              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
              `updated_at` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `permissions_actions` (`id`, `name`, `key`, `description`, `created_at`, `updated_at`) VALUES
            (1, 'Execute/Cancel Pending Transfer Requests', 'execute_cancel_pending_transfer_requests', NULL, '2018-11-05 06:32:33', NULL),
            (2, 'Initiate/Execute User Transfers', 'initiate_execute_user_transfers', NULL, '2018-11-05 06:32:36', NULL),
            (3, 'Send/Reply Internal Messages', 'send_reply_internal_messages', NULL, '2018-11-05 06:32:37', NULL),
            (4, 'Create profile', 'create_profile', NULL, '2018-08-24 14:32:02', NULL),
            (5, 'Approve/Cancel Registration Requests', 'approve_cancel_registration_requests', NULL, '2018-11-05 06:32:34', NULL),
            (6, 'Modify User Profiles', 'modify_user_profiles', NULL, '2018-08-24 14:32:02', NULL),
            (7, 'Modify Admin Profiles', 'modify_admin_profiles', NULL, '2018-08-24 14:32:02', NULL),
            (8, 'View User Profiles', 'view_user_profiles', NULL, '2018-11-05 06:32:35', NULL),
            (9, 'Create Accounts', 'create_accounts', NULL, '2018-08-24 14:32:02', NULL),
            (10, 'Create Accounts with Initial Balance', 'create_accounts_with_initial_balance', NULL, '2018-11-05 06:32:34', NULL),
            (12, 'Modify Accounts', 'modify_accounts', NULL, '2018-08-24 14:32:02', NULL),
            (13, 'Modify Account Types', 'modify_account_types', NULL, '2018-11-05 06:32:34', NULL),
            (14, 'Modify Cards', 'modify_cards', NULL, '2018-08-24 14:32:02', NULL),
            (15, 'Manual Debit/Credit Accounts', 'manual_debit_credit_accounts', NULL, '2018-11-05 06:32:35', NULL),
            (16, 'Import Transfer Request Updates', 'import_transfer_request_updates', NULL, '2018-08-24 14:32:02', NULL),
            (17, 'Manage Revenue', 'manage_revenue', NULL, '2018-11-05 06:32:39', NULL),
            (18, 'View User Reports', 'view_user_reports', NULL, '2018-11-05 06:32:38', NULL),
            (19, 'View General System Reports', 'view_general_system_reports', NULL, '2018-11-05 06:32:38', NULL),
            (20, 'View/Modify Settings', 'view_modify_settings', NULL, '2018-11-05 06:32:37', NULL),
            (22, 'Create/Modify IWT Bank Accounts', 'create_modify_iwt_bank_accounts', NULL, '2018-08-24 14:32:02', NULL),
            (25, 'View System Log', 'view_system_log', NULL, '2018-10-05 11:50:08', NULL),
            (27, 'Create/View Cards', 'create_view_cards', NULL, '2018-11-05 06:32:39', NULL),
            (28, 'Create/View Admin Profiles', 'create_view_admin_profiles', NULL, '2018-11-05 06:32:40', NULL),
            (29, 'View Revenue', 'view_revenue', NULL, '2018-11-05 06:34:12', NULL);

            CREATE TABLE `permissions_groups` (
              `id` int(10) UNSIGNED NOT NULL,
              `name` varchar(64) NOT NULL,
              `description` varchar(255) DEFAULT NULL,
              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
              `updated_at` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `permissions_groups` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
            (1, 'Administrator', 'Head Admin - Full Permission', NOW(), NOW());

            CREATE TABLE `permissions_groups_actions` (
              `group_id` int(10) UNSIGNED NOT NULL,
              `action_id` int(10) UNSIGNED NOT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `permissions_groups_actions` (`group_id`, `action_id`) VALUES
            (1, 1),
            (1, 2),
            (1, 3),
            (1, 4),
            (1, 5),
            (1, 6),
            (1, 7),
            (1, 8),
            (1, 9),
            (1, 10),
            (1, 12),
            (1, 13),
            (1, 14),
            (1, 15),
            (1, 16),
            (1, 17),
            (1, 18),
            (1, 19),
            (1, 20),
            (1, 22),
            (1, 25),
            (1, 27),
            (1, 28),
            (1, 29);

            CREATE TABLE `schema_migrations` (
              `version` bigint(20) NOT NULL,
              `dirty` tinyint(1) NOT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES
            (20181207174723, 0);


            ALTER TABLE `permissions_actions`
              ADD PRIMARY KEY (`id`),
              ADD UNIQUE KEY `name_UNIQUE` (`name`),
              ADD UNIQUE KEY `key_UNIQUE` (`key`);

            ALTER TABLE `permissions_groups`
              ADD PRIMARY KEY (`id`),
              ADD UNIQUE KEY `name_UNIQUE` (`name`);

            ALTER TABLE `permissions_groups_actions`
              ADD UNIQUE KEY `group_id_UNIQUE` (`group_id`,`action_id`),
              ADD KEY `FK_PERMISSIONS_GROUPS_ACTIONS_TO_GROUPS_idx` (`group_id`),
              ADD KEY `FK_PERMISSIONS_GROUPS_ACTIONS_TO_ACTIONS_idx` (`action_id`);

            ALTER TABLE `schema_migrations`
              ADD PRIMARY KEY (`version`);


            ALTER TABLE `permissions_actions`
              MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=30;

            ALTER TABLE `permissions_groups`
              MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;


            ALTER TABLE `permissions_groups_actions`
              ADD CONSTRAINT `FK_PERMISSIONS_GROUPS_ACTIONS_TO_ACTIONS` FOREIGN KEY (`action_id`) REFERENCES `permissions_actions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
              ADD CONSTRAINT `FK_PERMISSIONS_GROUPS_ACTIONS_TO_GROUPS` FOREIGN KEY (`group_id`) REFERENCES `permissions_groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
SQL;
    }
}
