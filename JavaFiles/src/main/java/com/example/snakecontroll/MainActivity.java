package com.example.snakecontroll;

import androidx.appcompat.app.AppCompatActivity;

import android.annotation.SuppressLint;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;

@SuppressLint("SetTextI18n")
public class MainActivity extends AppCompatActivity {
    Thread Thread1 = null;
    EditText etIP, etPort;
    String SERVER_IP;
    int SERVER_PORT;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        etIP = findViewById(R.id.etIP);
        etPort = findViewById(R.id.etPort);

        LoadIpAndPort();

        Button btnConnect = findViewById(R.id.btnConnect);
        btnConnect.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                SERVER_IP = etIP.getText().toString().trim();
                String portText = etPort.getText().toString().trim();
                if(SERVER_IP.length() > 0 && portText.length() > 0) {
                    SERVER_PORT = Integer.parseInt(portText);
                    Thread1 = new Thread(new Thread1());
                    Thread1.start();

                    try {
                        Thread.sleep(100);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                    if(output != null) {
                        SaveIpAndPort(SERVER_IP, SERVER_PORT);
                        startControll();
                    }else{
                        System.out.println("Error connecting");
                    }
                }
            }
        });
    }
    public void SaveIpAndPort(String ip, int port) {
        SharedPreferences sharedPreferences= this.getSharedPreferences("serverConn", Context.MODE_PRIVATE);
        SharedPreferences.Editor editor = sharedPreferences.edit();
        editor.putString("IP", ip);
        editor.putInt("PORT", port);
        editor.apply();
    }
    public void LoadIpAndPort() {
        SharedPreferences sharedPreferences= this.getSharedPreferences("serverConn", Context.MODE_PRIVATE);
        if(sharedPreferences != null) {
            SERVER_IP = sharedPreferences.getString("IP", "");
            SERVER_PORT = sharedPreferences.getInt("PORT",8080);
            etIP.setText(SERVER_IP);
            System.out.println(SERVER_PORT);
            etPort.setText(Integer.toString(SERVER_PORT));
        }
    }

    public void startControll() {
        Intent intent = new Intent(this, ControllActivity.class);
        startActivity(intent);
    }

    public static PrintWriter output;
    public BufferedReader input;
    class Thread1 implements Runnable {
        @Override
        public void run() {
            Socket socket;
            try {
                socket = new Socket(SERVER_IP, SERVER_PORT);
                output = new PrintWriter(socket.getOutputStream());
                input = new BufferedReader(new InputStreamReader(socket.getInputStream()));
                System.out.println("Connected\n");
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}