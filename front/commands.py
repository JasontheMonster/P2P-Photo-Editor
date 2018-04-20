import ps
import front

from Tkinter import *
import tkFileDialog
import tkSimpleDialog
import cv2

img = None

def selectImg(btn):
	path = tkFileDialog.askopenfilename()
	if len(path) > 0:
		btn.destroy()
		global img
		img = cv2.imread(path)
		front.showImg()
		return True
	return False

def invite():
	IP = tkSimpleDialog.askstring("IP address","IP: ")
	print IP

def bright():
	global img
	img = ps.brighten(img, 25)
	front.showImg()

def saturate():
	global img
	img = ps.saturation(img, 15)
	front.showImg()

def blur():
	global img
	size = img.shape
	klen = min(size[0],size[1])
	img = ps.bilareralBlur(img, klen/10, klen/10)
	front.showImg()

def rotate():
	global img
	img = ps.rotateImg(img, 180)
	front.showImg()

def contrast():
	global img
	img = ps.contrast(img,100)
	front.showImg()