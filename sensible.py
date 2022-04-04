#! /usr/bin/env python3
#-*- coding:utf-8 -*-
import sys; sys.dont_write_bytecode = True;

import argparse
import itertools
import json
import math
import os
import subprocess
import textwrap

import curses
import curses.panel
import curses.textpad

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
    self.run_plays = False
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
    if self.run_plays:
      self.run_playbooks()


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

  def run_playbooks(self):
    for option in self.options:
      if option['selected']:
        playbook_path = option['path']
        ansible_cmd = f"ansible-playbook {playbook_path}"
        os.system(ansible_cmd)


  ############
  #
  def notify(self, key, value):
    setattr(self, key, value)
    return value

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
  def slice_text(self, text, width, padding):
    sliced = []
    if len(text) < (width + padding):
      return [text]
    line = ""
    parts = list(itertools.chain.from_iterable(zip(text.split(), itertools.repeat(' '))))[:-1]
    for part in parts:
      if len(line) + len(part) < width:
        line += part
      else:
        sliced.append(line)
        line = part
      if part == parts[-1]:
        sliced.append(line)
    #sliced.append('...')
    return sliced

  def create_window(self, h, w, x, y ):
    window = curses.newwin( h, w, x, y )
    window.erase()
    window.immedok(True)
    window.box()
    window.bkgd(" ", curses.color_pair(4))
    window.refresh( )
    return window

  def render_title(self):
    text = self.elements['title']
    height, width = self.stdscr.getmaxyx()
    start_y = int((width // 2) - (len(text) // 2) - len(text) % 2)
    self.stdscr.addstr(0, 0, f"{' ' * self.get_width()}", curses.color_pair(3) )
    self.stdscr.addstr(0, start_y, text, curses.color_pair(3))

  def render_chyron(self):
    text = " | ".join(f'{k}: {v}' for k,v in self.elements['chyron'].items())
    height, width = self.stdscr.getmaxyx()
    self.stdscr.addstr(height-1, 0, " " * (width-1), curses.color_pair(3))
    self.stdscr.addstr(height-1, 0, text, curses.color_pair(3))

  def render_left_panel(self):
    max_x = math.floor(((self.get_width() / 9 )) * 6 )
    max_y = math.floor((self.get_height() -3))
    window = self.create_window( max_y, max_x, 1, 0 )
    for i, option in enumerate(self.options):
      highlight = curses.color_pair(4)
      # Check if selection is a seperator
      # if 'seperator' in option['tags']:
      #   padding = (max_x - ( len(option['name']) + 3 ))
      #   window.addstr((i + 2), 2, f"> {'-' * padding}", curses.color_pair(1))
      if option['selected']:
        highlight = curses.color_pair(5)
      if i == self.position:
        if not option['selected']:
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
    window = self.create_window( (self.get_height() -3), max_x, 1, _x )
    cur_selection = self.options[self.position]
    content = [
      f"Name: {cur_selection['name']}",
      f"Description:"
    ]
    content += self.slice_text(cur_selection['description'], max_x, 2)
    for i, line in enumerate(content):
      # if len(line) <= max_x - 2:
      #   window.addstr(i + 1, 2, line, curses.color_pair(1))
      window.addnstr((i + 2), 2, textwrap.fill(f"{line}", (max_x -2)), curses.color_pair(1))
    panel = curses.panel.new_panel(window)
    return window, panel


  ############
  # TUI
  def run(self, stdscr):
    self.stdscr = stdscr
    self.stdscr.border( 0 )
    self.stdscr.keypad(1)

    curses.noecho()
    curses.cbreak()
    curses.curs_set( 0 )

    curses.start_color()
    curses.init_pair(4, curses.COLOR_WHITE, curses.COLOR_BLACK)
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)
    curses.init_pair(5, curses.COLOR_GREEN, curses.COLOR_BLACK)
    curses.init_pair(6, curses.COLOR_BLUE, curses.COLOR_BLACK)
    curses.init_pair(9, curses.COLOR_WHITE, curses.COLOR_BLUE)

    k = 0
    cursor_y = 0
    cursor_y_max = len(self.options) -1

    while (k != ord('q')):
      # Initialization
      # Clear and refresh the screen for a blank canvas
      stdscr.clear()
      stdscr.refresh()

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
            not self.options[self.position]['selected'] )
      elif k == curses.KEY_ENTER or k == 10:
        self.run_plays = True
        break

      self.position = cursor_y

      # Render title
      self.render_title()
      # Render status bar
      self.render_chyron()
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


def main(**kwargs):
  app = Sensible(**kwargs)


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
