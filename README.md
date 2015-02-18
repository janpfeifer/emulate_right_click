emulate_right_click translates CONTROL+left_click to a right_click mouse event.

I don't like the double finger tap on my touchpad (trackpad) so I tried this instead. This
can be useful if you mouse doesn't have a right button -- or if it is broken.

Disclaimer: this is my first windows program ... likely it doesn't work as intented in
every case -- although it's working for me know. Use at your own risk. Check the source
code if you are worried, it is small.

The combination CONTROL+left_click is hardcoded for now, but changing it is relatively
simple.

Requirements:

  * Go compiler: http://golang.org/
  * gcc compiler for windows:
    http://mingw-w64.sourceforge.net/
  * Golang package, installed with
    go get github.com/AllenDang/w32
    go install github.com/AllenDang/w32

Documentation reference:
  * godoc for package github.com/AllenDang/w32, and looking at the source code.
  * Example code in http://play.golang.org/p/kwfYDhhiqk
  * C++ example of mouse emulation (oudated but it was helpful to get an idea):
    http://www.codeproject.com/Articles/194265/Mouse-emulating-software
  * Windows API Index:
    https://msdn.microsoft.com/en-US/library/windows/desktop/ff818516(v=vs.85).aspx
