from Tkinter import *
from PIL import Image
from PIL import ImageTk
from threading import Thread
import tkFileDialog
import tkSimpleDialog
import cv2
import socket
import ps
import os
import sys

class App:
	def __init__(self,master):
		self.isAlive = True
		self.master = master
		self.panel = None
		self.img = None
		self.listeningAddr = (sys.argv[1].split(':')[0], int(sys.argv[1].split(':')[1]))
		self.sendingAddr = (sys.argv[2].split(':')[0], int(sys.argv[2].split(':')[1]))

		self.image_source = ""
		self.Btn_invite = None
		self.Btn_bright = None
		self.Btn_saturate = None
		self.Btn_blur = None
		self.Btn_rotate = None
		self.Btn_contrast = None
		self.Btn_selectImg = Button(self.master, text='Select image', width = 50, height=30, bg='gold', command=lambda: self.__start())
		self.Btn_selectImg.grid(row=0,column=0)
		self.Btn_quit = Button(self.master, text='quit', width = 50, height=2, bg='gold', command=lambda: self.__quit())
		self.Btn_quit.grid(row=6,column=0)

		self.ops = {'blur': self.__blur, 'bright': self.__bright, 'saturate': self.__saturate, 'rotate': self.__rotate, 'contrast': self.__contrast}
		Thread(target=self.__listener).start()

	def __listener(self):
		sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
		sock.bind(self.listeningAddr)
		sock.listen(1)
		while self.isAlive:
			conn,_ = sock.accept()
			logEty = conn.recv(1024)
			if logEty == 'quit':
				break
			self.__eventHandler(logEty)
			conn.close()
		sock.close()

	def __eventHandler(self, logEty):
		if logEty.split('@')[0] == "Image":
				print "receive image from others!!!!"
				image_path = logEty.split('@')[1].replace("\\","/")
				self.img = cv2.imread(image_path)
				self.__shrinkImg()
				self.Btn_selectImg.destroy()
				self.__showImg()
				self.__preparePanel()
				self.image_source = image_path
				self.__send("PATH:"+self.image_source)
		else:
			print "***"+logEty+"***"
			if logEty != 'quit':
				self.ops[logEty]()

	def __send(self, msg):
		sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
		sock.connect(self.sendingAddr)
		sock.send(msg)
		sock.close()

	def __quit(self):
		self.master.destroy() # close tne window
		self.isAlive = False
		self.__send('quit')

	def __selectImg(self):
		path = tkFileDialog.askopenfilename()
		print path
		if len(path) > 0:
			self.Btn_selectImg.destroy()
			self.img = cv2.imread(path)
			self.__shrinkImg()
			self.image_source = path
			self.__send("PATH:"+self.image_source)
			self.__showImg()
			return True
		return False

	def __preparePanel(self):
		self.Btn_invite = Button(self.master, text='invite', height=6, width=15, bg='gold', command=lambda: self.__invite())
		self.Btn_invite.grid(row=0,column=1)

		self.Btn_bright = Button(self.master, text='bright', height=6, width=15, bg='gold', command=lambda: self.__send('bright'))
		self.Btn_bright.grid(row=1,column=1)

		self.Btn_saturate = Button(self.master, text='saturate', height=6, width=15, bg='gold', command=lambda: self.__send('saturate'))
		self.Btn_saturate.grid(row=2,column=1)

		self.Btn_blur = Button(self.master, text='blur', height=6, width=15, bg='gold', command=lambda: self.__send('blur'))
		self.Btn_blur.grid(row=3,column=1)

		self.Btn_rotate = Button(self.master, text='rotate', height=6, width=15, bg='gold', command=lambda: self.__send('rotate'))
		self.Btn_rotate.grid(row=4,column=1)

		self.Btn_contrast = Button(self.master, text='contrast', height=6, width=15, bg='gold', command=lambda: self.__send('contrast'))
		self.Btn_contrast.grid(row=5,column=1)

	def __start(self):
		flag = self.__selectImg()
		if flag:
			self.__preparePanel()
			
	def __shrinkImg(self):
		h,_ = self.img.shape[:2]
		max_height = 600
		# only shrink if img is bigger than required
		if max_height < h:
		    # get scaling factor
		    scaling_factor = max_height / float(h)
		    # resize image
		    self.img = cv2.resize(self.img, None, fx=scaling_factor, fy=scaling_factor, interpolation=cv2.INTER_AREA)


	def __showImg(self):
		displayImg = cv2.cvtColor(self.img, cv2.COLOR_BGR2RGB)
		displayImg = Image.fromarray(displayImg)
		displayImg = ImageTk.PhotoImage(displayImg)
		if self.panel is None:
			self.panel = Label(image=displayImg)
			self.panel.image = displayImg
			self.panel.grid(row=0,column=0,rowspan=6)
		else:
			self.panel.configure(image=displayImg)
			self.panel.image = displayImg

	def __invite(self):
		IP = tkSimpleDialog.askstring("IP address","IP: ")
		self.__send("invite"+IP)

	def __bright(self):
		self.img = ps.brighten(self.img, 25)
		self.__showImg()

	def __saturate(self):
		self.img = ps.saturation(self.img, 15)
		self.__showImg()

	def __blur(self):
		size = self.img.shape
		klen = min(size[0],size[1])
		self.img = ps.bilareralBlur(self.img, klen/10, klen/10)
		self.__showImg()

	def __rotate(self):
		self.img = ps.rotateImg(self.img, 180)
		self.__showImg()

	def __contrast(self):
		self.img = ps.contrast(self.img,100)
		self.__showImg()


root = Tk()
root.title("Collaborative Photo Editor")
root.config(background='DarkOrchid1')
app = App(root)
root.mainloop()
