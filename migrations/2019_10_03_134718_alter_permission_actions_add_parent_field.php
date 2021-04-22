<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AlterPermissionActionsAddParentField extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('permissions_actions', function (Blueprint $table) {
            $table->dropForeign('parent_id');
            $table->dropForeign('category_id');
            $table->dropColumn('parent_id');
            $table->dropColumn('category_id');
            $table->dropColumn('is_hidden');
            $table->dropColumn('sort');
        });
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('permissions_actions', function (Blueprint $table) {
            $table->integer('sort')->default(0);
            $table->boolean('is_hidden')->default(false);
            $table->unsignedInteger('parent_id')->nullable();
            $table->unsignedInteger('category_id')->nullable();

            $table->foreign('category_id')->references('id')->on('permissions_categories')
                ->onUpdate('cascade')->onDelete('set null');

            $table->foreign('parent_id')->references('id')->on('permissions_actions')
                ->onUpdate('cascade')->onDelete('cascade');
        });

        $this->insertNewPermissions();
    }

    private function insertNewPermissions()
    {
        $permissions = [
            'view_accounts' => [
                'name' => 'View Accounts',
            ],
        ];

        foreach ($permissions as $key => $permission) {
            DB::insert('INSERT INTO permissions_actions (`name`, `key`) VALUES (?, ?)',
                [$permission['name'], $key]);
        }
    }
}
