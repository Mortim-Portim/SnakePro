package com.example.snakecontroll;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.view.MotionEventCompat;

import android.os.Bundle;
import android.util.DisplayMetrics;
import android.view.MotionEvent;
import android.view.View;
import android.widget.ImageView;

/**
 * Sends Data to MainActivity.output
 * Format:
 * "[left][right][up][down][button]"
 */
public class ControllActivity extends AppCompatActivity {
    public float width, height;
    public int leftPosX, leftPosY, downXPos, downYPos;
    public boolean left, right, up, down, button = false;
    public DisplayMetrics displayMetrics;
    public ImageView moveBtn;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_controll);
        hideSystemUI();

        moveBtn = findViewById(R.id.moveBtn);

        displayMetrics = new DisplayMetrics();
    }

    @Override
    public boolean onTouchEvent(MotionEvent event) {
        String oldState = ButtonState();
        int action = MotionEventCompat.getActionMasked(event);
        int pointers = event.getPointerCount();
        this.getWindowManager().getDefaultDisplay().getMetrics(displayMetrics);
        width = displayMetrics.widthPixels;
        height = displayMetrics.heightPixels;

        for(int i=0; i<pointers; i++) {
            int xPos = (int) MotionEventCompat.getX(event, i);
            int yPos = (int) MotionEventCompat.getY(event, i);

            if(action == MotionEvent.ACTION_DOWN || action == MotionEvent.ACTION_POINTER_DOWN) {
                if(xPos > width/4 && xPos < width/4*3 && yPos > height/4 && yPos < height/4*3) {
                    hideSystemUI();
                }
            }

            if(xPos > width/2) {
                switch (action) {
                    case MotionEvent.ACTION_DOWN:
                    case MotionEvent.ACTION_POINTER_DOWN:
                        if(i == event.getActionIndex()) {
                            button = true;
                        }
                        break;
                    case MotionEvent.ACTION_UP:
                    case MotionEvent.ACTION_POINTER_UP:
                        if(i == event.getActionIndex()) {
                            button = false;
                        }
                        break;
                }
            }else{
                switch (action) {
                    case MotionEvent.ACTION_MOVE:
                        ToggleDir(false);
                        SetPointer(xPos,yPos);
                        ToggleDir(true);
                        break;
                    case MotionEvent.ACTION_DOWN:
                    case MotionEvent.ACTION_POINTER_DOWN:
                        if(i == event.getActionIndex()) {
                            downXPos = xPos;
                            downYPos = yPos;
                            moveBtn.setX(downXPos-moveBtn.getWidth()/2);
                            moveBtn.setY(downYPos-moveBtn.getHeight()/2);
                        }
                        break;
                    case MotionEvent.ACTION_UP:
                    case MotionEvent.ACTION_POINTER_UP:
                        if(i == event.getActionIndex()) {
                            ResetPointer();
                        }
                        break;
                }
            }

        }
        String newState = ButtonState();
        if(!oldState.equalsIgnoreCase(newState)) {
            Send(newState);
        }
        return super.onTouchEvent(event);
    }
    public String ButtonState() {
        String out = "";
        out += BoolToStr(left);
        out += BoolToStr(right);
        out += BoolToStr(up);
        out += BoolToStr(down);
        out += BoolToStr(button);
        return out;
    }
    public static String BoolToStr(boolean b) {
        if(b) {
            return "1";
        }else {
            return "0";
        }
    }
    public void ResetPointer() {
        left = false; right = false; up = false; down = false;
    }
    public void SetPointer(int x, int y) {
        leftPosX = x;
        leftPosY = y;
    }
    public void ToggleDir(boolean val) {
        if(Math.abs(leftPosX-downXPos) > Math.abs(leftPosY-downYPos)) {
            if(Math.abs(leftPosX-downXPos) > width/50) {
                if (leftPosX - downXPos < 0) {
                    left = val;
                } else {
                    right = val;
                }
            }
        }else{
            if(Math.abs(leftPosY-downYPos) > width/50) {
                if (leftPosY - downYPos < 0) {
                    up = val;
                } else {
                    down = val;
                }
            }
        }
    }

    @Override
    protected void onResume() {
        super.onResume();
        hideSystemUI();
    }

    private void hideSystemUI() {
        View decorView = getWindow().getDecorView();
        decorView.setSystemUiVisibility(
                View.SYSTEM_UI_FLAG_IMMERSIVE
                        | View.SYSTEM_UI_FLAG_LAYOUT_STABLE
                        | View.SYSTEM_UI_FLAG_LAYOUT_HIDE_NAVIGATION
                        | View.SYSTEM_UI_FLAG_LAYOUT_FULLSCREEN
                        | View.SYSTEM_UI_FLAG_HIDE_NAVIGATION
                        | View.SYSTEM_UI_FLAG_FULLSCREEN);
    }

    public void Send(String msg) {
        new Thread(new Thread3(msg)).start();
    }
    class Thread3 implements Runnable {
        private String message;
        Thread3(String message) {
            this.message = message;
        }
        @Override
        public void run() {
            MainActivity.output.write(message+"\n");
            MainActivity.output.flush();
        }
    }
}