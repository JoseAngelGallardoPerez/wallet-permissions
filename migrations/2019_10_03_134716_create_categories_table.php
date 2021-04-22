<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCategoriesTable extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('permissions_categories');
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('permissions_categories', function (Blueprint $table) {
            $table->increments('id');
            $table->string('name', 255);
            $table->integer('sort');
        });

        $this->insertCategories();
    }

    private function insertCategories()
    {
        $categories = [
            [
                'name' => 'Requests',
                'sort' => 1000,
            ],
            [
                'name' => 'Accounts',
                'sort' => 950,
            ],
            [
                'name' => 'Profiles',
                'sort' => 900,
            ],
            [
                'name' => 'Transfers',
                'sort' => 850,
            ],
            [
                'name' => 'Messages',
                'sort' => 800,
            ],
            [
                'name' => 'System Log',
                'sort' => 700,
            ],
            [
                'name' => 'Settings',
                'sort' => 650,
            ],
            [
                'name' => 'Reports',
                'sort' => 600,
            ],
        ];

        foreach ($categories as $key => $category) {
            DB::insert('INSERT INTO permissions_categories (`name`, `sort`) VALUES (?, ?)',
                [$category['name'], $category['sort']]);
        }
    }
}
