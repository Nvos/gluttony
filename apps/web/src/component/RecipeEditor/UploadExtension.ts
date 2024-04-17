import {Extension} from '@tiptap/core';
import { Plugin, PluginKey } from '@tiptap/pm/state'

const ImageDnd = Extension.create({
  name: 'dnd',
  addProseMirrorPlugins() {
    return [
      new Plugin({
        key: new PluginKey('eventHandler'),
        props: {
          handleDOMEvents: {
            drop: (view, event) => {

            }
          }
        }
      })
    ]
  }
})