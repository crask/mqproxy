<?php
    $ch =curl_init("http://127.0.0.1:9090/produce?format=json");
    curl_setopt($ch, CURLOPT_HEADER, 0);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_BINARYTRANSFER, true);

    $data = array(
        'topic' => 'test',
        'partitionKey' => '123',
        'data' => array(
            'uid'     => 123,
            'uname'   => 'crask',
            'content' => 'welcome to crask',
        ),
    );

    $json_str = json_encode($data);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $json_str);
    $res = curl_exec($ch);
    var_dump($res);
?>
