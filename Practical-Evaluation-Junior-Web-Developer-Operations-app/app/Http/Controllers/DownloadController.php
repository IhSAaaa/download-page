<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\Download;
use Illuminate\Support\Facades\Mail;
use App\Mail\DownloadConfirmation;

class DownloadController extends Controller
{
    public function store(Request $request)
    {
        $data = $request->validate([
            'firstname' => 'required|string',
            'surname' => 'required|string',
            'company' => 'required|string',
            'country' => 'required|string',
            'prefix' => 'required|string',
            'phone' => 'required|string',
            'email' => 'required|email',
        ]);
        Download::create($data);

        $datas = array(
            'firstname' =>  $request->firstname,
            'surname'   =>  $request->surname,
            'company'   =>  $request->company,
            'country'   =>  $request->country,
            'prefix'    =>  $request->prefix,
            'phone'     =>  $request->phone,
            'email'     =>  $request->email
        );

        Mail::to($request->email)->send(new DownloadConfirmation($datas));

        $file = public_path()."/TKI-Membership-Brochure.pdf";
        $headers = array('Content-Type: application/pdf',);
        return response()->download($file, 'TKI-Membership-Brochure.pdf', $headers);
        
    }
}