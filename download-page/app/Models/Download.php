<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Download extends Model
{
    protected $fillable = [
        'firstname',
        'surname',
        'company',
        'country',
        'prefix',
        'phone',
        'email',
    ];
}