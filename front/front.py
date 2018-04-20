from Tkinter import *
from PIL import Image
from PIL import ImageTk

import cv2

import commands

def showImg():
	global panel
	displayImg = cv2.cvtColor(commands.img, cv2.COLOR_BGR2RGB)
	displayImg = Image.fromarray(displayImg)
	displayImg = ImageTk.PhotoImage(displayImg)
	if panel is None:
		panel = Label(image=displayImg)
		panel.image = displayImg
		panel.grid(row=0,column=0,rowspan=10)
	else:
		panel.configure(image=displayImg)
		panel.image = displayImg

def init(btn):
	flag = commands.selectImg(btn)
	if flag:

		global Btn_invite
		Btn_invite = Button(root, text='invite', height=5, width=15, bg='indian red', command=lambda: commands.invite())
		Btn_invite.grid(row=0,column=1)

		global Btn_bright
		Btn_bright = Button(root, text='bright', height=5, width=15, bg='indian red', command=lambda: commands.bright())
		Btn_bright.grid(row=1,column=1)

		global Btn_saturate
		Btn_saturate = Button(root, text='saturate', height=5, width=15, bg='indian red', command=lambda: commands.saturate())
		Btn_saturate.grid(row=2,column=1)

		global Btn_blur
		Btn_blur = Button(root, text='blur', height=5, width=15, bg='indian red', command=lambda: commands.blur())
		Btn_blur.grid(row=3,column=1)

		global Btn_rotate
		Btn_rotate = Button(root, text='rotate', height=5, width=15, bg='indian red', command=lambda: commands.rotate())
		Btn_rotate.grid(row=4,column=1)

		global Btn_contrast
		Btn_contrast = Button(root, text='contrast', height=5, width=15, bg='indian red', command=lambda: commands.contrast())
		Btn_contrast.grid(row=5,column=1)

panel = None
Btn_selectImg = None
Btn_invite = None
Btn_bright = None
Btn_saturate = None
Btn_blur = None
Btn_rotate = None
Btn_contrast = None

if __name__ == "__main__":
	root = Tk()
	root.title("Collaborative Photo Editor")
	root.config(background='light green')
	Btn_selectImg = Button(root, text='Select an image', width = 50, height=30, bg='indian red', command=lambda: init(Btn_selectImg))
	Btn_selectImg.grid(row=0,column=0)
	root.mainloop()