#! /usr/bin/env python3
#-*- coding:utf-8 -*-
import sys; sys.dont_write_bytecode = True;

import argparse
import json
import math
import os

import curses
import curses.panel

from pathlib import Path
from pprint import pprint

import yaml


__app__ = "Sensible"
__description__ = "Ansible Playbook TUI"


#########################################
class Sensible:
  subscriptions = {
    'dir': {
      'validate': 'check_path',
      'notify':   'find_playbooks'
      #'attach'
      #'update'
      #'render'
    }
  }

  def __init__(self,**kwargs):
    self.options = { }
    self.position = 0
    self.elements = {
      'title':  "Sensible - Ansible Playbook TUI",
      'chyron': {
        'space': "Select playbooks",
        'enter': "Run playbooks",
        'q': "Quit",
        'a': "About"
      }
    }


    for key, value in kwargs.items(): self.attach(key, value)

    curses.wrapper(self.run)



  ############
  # Utils
  def attach(self, key, value):
    if key in self.subscriptions.keys():
      observer = self.subscriptions[key]
      observer_keys = observer.keys()
      if 'validate' in observer_keys:
        validate = getattr(self,  observer['validate'])
        if validate(value):
          setattr(self, key, value)
      if 'notify' in observer_keys:
        notify = getattr(self,  observer['notify'])
        notify(value)
        #setattr(self, key, value)
    else:
      raise Exception("[!] Invalid argument")

  def notify(self, key, value):
    setattr(self, key, value)
    return value

  def check_path(self, path):
    if not Path(path).is_dir():
      print("[!] Directory does not exist")
      sys.exit(1)
    return True


  ############
  # Playbooks
  def parse_playbook(self, playbook_path):
    start_pattern = '### sensible ###'
    end_pattern   = '### /sensible ###'
    content = Path(playbook_path).read_text()
    if start_pattern in content and end_pattern in content:
      header_str = content.partition(end_pattern)[0].rpartition(start_pattern)[2]
      header_str = header_str.translate({ord(c): None for c in '!@#$'})
      header_yaml = yaml.load(header_str, Loader=yaml.SafeLoader)
      header = dict((key, val) for k in header_yaml for key, val in k.items())
      return header
    return False


  def find_playbooks(self, playbook_dir):
    playbooks = [None] * 50

    files = [ f for f in Path(playbook_dir).iterdir()
      if f.is_file() and f.suffix in ['.yml', '.yaml'] ]

    for f in files:
      parsed = self.parse_playbook(f)
      if parsed:
        parsed['path'] = str(f)
        parsed['selected'] = False
        if type(parsed['index']) == int:
          playbooks[parsed['index']] = parsed
        else:
          playbooks.append(parsed)

    playbooks = list(filter(None, playbooks))

    if not playbooks:
      print("[!] No playbooks found")
      sys.exit(1)
    self.options = playbooks

  ############
  #
  def get_height( self ):
    height = int( os.popen('stty size').read( ).split( )[0].strip( ) )
    # return self.notify('height', height, self.window)
    return height

  def get_width( self ):
    width = int( os.popen('stty size').read( ).split( )[1].strip( ) )
    # return self.notify('width', width, self.window)
    return width

  def center_text(self, text, width):
    return int((width // 2) - (len(text) // 2) - len(text) % 2)

  ############
  #
  def render_title(self, text):
    text = "Sensible - Ansible Playbook TUI"
    height, width = self.stdscr.getmaxyx()
    start_y = int((width // 2) - (len(text) // 2) - len(text) % 2)
    max_y = (self.get_width( ))
    self.stdscr.addstr(0, 0, f"{' ' * max_y}", curses.color_pair(3) )
    self.stdscr.addstr(0, start_y, text, curses.color_pair(3))


  def render_chyron(self):
    #self.elements['chyron']
    text = "test"
    height, width = self.stdscr.getmaxyx()
    self.stdscr.attron(curses.color_pair(3))
    self.stdscr.addstr(height-1, 0, text)
    self.stdscr.addstr(height-1, len(text), " " * (width - len(text) - 1))
    self.stdscr.attroff(curses.color_pair(3))

  def render_left_panel(self):
    max_x = math.floor(((self.get_width() / 9 )) * 6 )
    max_y = math.floor((self.get_height() -2))
    # curses.newwin( h, w, y, x )
    window = curses.newwin( max_y, max_x, 1, 0 )
    window.erase()
    window.box()
    window.bkgd(" ", curses.color_pair(4))
    window.refresh( )
    for i, option in enumerate(self.options):
      highlight = curses.color_pair(4)
      # Check if selection is a seperator
      # if 'seperator' in option['tags']:
      #   padding = (max_x - ( len(option['name']) + 3 ))
      #   window.addstr((i + 2), 2, f"> {'-' * padding}", curses.color_pair(1))
      if option['selected']:
        highlight = curses.color_pair(5)
      if i == self.position:
        highlight = curses.color_pair(1)
        window.addstr((i + 2), 2, f"> {option['name']}", highlight)
      else:
        window.addstr((i + 2), 2, f"  {option['name']}", highlight)
      i=i+1
    panel = curses.panel.new_panel(window)
    return window, panel

  def render_right_panel(self):
    _x = math.floor(((self.get_width() / 9 )) * 6 )
    max_x = math.floor(((self.get_width() / 9 )) * 3 )
    max_y = math.floor((self.get_height() -2))
    # curses.newwin( h, w, y, x )
    window = curses.newwin( max_y, max_x, 1, _x )
    window.erase()
    window.box()
    window.bkgd(" ", curses.color_pair(4))
    window.refresh( )
    cur_selection = self.options[self.position]
    #meta = json.dumps(cur_selection, sort_keys=False, indent=2)
    #meta = yaml.dump(cur_selection, sort_keys=False, default_flow_style=False)
    meta = f"{cur_selection['name']}\n\n{cur_selection['description']}"
    window.addstr((0 + 2), 2, f"{meta}", curses.color_pair(1))
    panel = curses.panel.new_panel(window)
    return window, panel

  ############
  # TUI
  def run(self, stdscr):
    self.stdscr = stdscr
    # Clear and refresh the screen for a blank canvas
    stdscr.clear()
    stdscr.refresh()
    ###
    curses.noecho()
    curses.cbreak()
    self.stdscr.border( 0 )
    curses.curs_set( 0 )
    self.stdscr.keypad(1)
    # Start colors in curses
    curses.start_color()
    curses.init_pair(4, curses.COLOR_WHITE, curses.COLOR_BLACK)
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)
    curses.init_pair(5, curses.COLOR_GREEN, curses.COLOR_BLACK)
    curses.init_pair(6, curses.COLOR_BLUE, curses.COLOR_BLACK)
    # curses.init_pair(2, curses.COLOR_GREEN, curses.COLOR_BLACK)
    # curses.init_pair(3, curses.COLOR_CYAN, curses.COLOR_BLACK)
    # curses.init_pair(4, curses.COLOR_YELLOW, curses.COLOR_BLACK)
    # curses.init_pair(5, curses.COLOR_GREEN, curses.COLOR_BLACK)
    # curses.init_pair(6, curses.COLOR_BLACK, curses.COLOR_WHITE)
    # curses.init_pair(7, curses.COLOR_RED, curses.COLOR_BLACK)
    # curses.init_pair(8, curses.COLOR_WHITE, curses.COLOR_WHITE)
    curses.init_pair(9, curses.COLOR_WHITE, curses.COLOR_BLUE)

    # Loop where k is the last character pressed
    k = 0
    cursor_x = 0
    cursor_y = 0



    while (k != ord('q')):
      # Initialization
      stdscr.clear()
      height, width = stdscr.getmaxyx()

      # Render title
      self.render_title(self.elements['title'])
      # Render status bar
      self.render_chyron()

      cursor_y_max = len(self.options) -1
      if k == curses.KEY_DOWN:
        if (cursor_y + 1) >= cursor_y_max:
          cursor_y = cursor_y_max
        else:
          cursor_y += 1
      elif k == curses.KEY_UP:
        if (cursor_y - 1) < 0:
          cursor_y = 0
        else:
          cursor_y -= 1

      elif k == ord(' '):
          self.options[self.position]['selected'] = (
            not self.options[self.position]['selected']
          )
      elif k == curses.KEY_ENTER:
          cursor_x = cursor_x - 1

      # cursor_y = max(0, len(self.options) - 1)
      # cursor_y = min(height-1, cursor_y)
      self.position = cursor_y

      # Declaration of strings
      title = "Curses example"[:width-1]
      subtitle = "Written by Clay McLeod"[:width-1]
      keystr = "Last key pressed: {}".format(k)[:width-1]
      if k == 0:
          keystr = "No key press detected..."[:width-1]

      # Centering calculations
      start_x_title = int((width // 2) - (len(title) // 2) - len(title) % 2)
      start_x_subtitle = int((width // 2) - (len(subtitle) // 2) - len(subtitle) % 2)
      start_x_keystr = int((width // 2) - (len(keystr) // 2) - len(keystr) % 2)
      start_y = int((height // 2) - 2)

      # # Rendering some text
      # whstr = "Width: {}, Height: {}".format(width, height)
      # stdscr.addstr(0, 0, whstr, curses.color_pair(1))
      # Header
      # whstr = "Sensible - Ansible Playbook TUI"
      # stdscr.addstr(0, 0, whstr, curses.color_pair(3))


      ## Window/Panel
      win1, panel1 = self.render_left_panel()
      curses.panel.update_panels(); stdscr.refresh()
      panel1.top(); curses.panel.update_panels(); stdscr.refresh()
      ## Window/Panel
      win1, panel2 = self.render_right_panel()
      curses.panel.update_panels(); stdscr.refresh()
      panel2.top(); curses.panel.update_panels(); stdscr.refresh()


      # Refresh the screen
      stdscr.refresh()
      # curses.flushinp()

      # Wait for next input
      k = stdscr.getch()

#########################################

def main(**kwargs):
  # print(kwargs)
  app = Sensible(**kwargs)


#########################################

if __name__ == "__main__":
  parser = argparse.ArgumentParser(
    prog=__app__,
    description=__description__
  )

  parser.add_argument(
    "-d",
    "--dir",
    nargs='?',
    required=True,
    help=""
  )

  args = vars( parser.parse_args() )

  main(**args)

