<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Mail;
use App\Mail\DownloadConfirmation;

class SendMail extends Controller
{
    function index()
    {
     return view('test');
    }

    function send(Request $request)
    {
     $this->validate($request, [
        'firstname' => 'required|string',
        'surname' => 'required|string',
        'company' => 'required|string',
        'country' => 'required|string',
        'prefix' => 'required|string',
        'phone' => 'required|string',
        'email' => 'required|email',
    ]);

        $data = array(
            'firstname' =>  $request->firstname,
            'surname'   =>  $request->surname,
            'company'   =>  $request->company,
            'country'   =>  $request->country,
            'prefix'    =>  $request->prefix,
            'phone'     =>  $request->phone,
            'email'     =>  $request->email
        );

     Mail::to($request->email)->send(new DownloadConfirmation($data));
     return back();
    }
}
