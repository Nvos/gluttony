import { Node, ReactRenderer } from '@tiptap/react';
import tippy, { Instance, Tippy } from 'tippy.js';
import { MentionOptions } from '@tiptap/extension-mention';

import { IngredientList } from './IngredientList';

export const suggestion: MentionOptions['suggestion'] = {
  items: ({ query }) => {
    return [
      'Lea Thompson',
      'Cyndi Lauper',
      'Tom Cruise',
      'Madonna',
      'Jerry Hall',
      'Joan Collins',
      'Winona Ryder',
      'Christina Applegate',
      'Alyssa Milano',
      'Molly Ringwald',
      'Ally Sheedy',
      'Debbie Harry',
      'Olivia Newton-John',
      'Elton John',
      'Michael J. Fox',
      'Axl Rose',
      'Emilio Estevez',
      'Ralph Macchio',
      'Rob Lowe',
      'Jennifer Grey',
      'Mickey Rourke',
      'John Cusack',
      'Matthew Broderick',
      'Justine Bateman',
      'Lisa Bonet',
    ]
      .filter((item) => item.toLowerCase().startsWith(query.toLowerCase()))
      .slice(0, 5);
  },

  render: () => {
    let component: ReactRenderer;
    let popup: Instance;

    return {
      onStart: (props) => {
        component = new ReactRenderer(IngredientList, {
          props,
          editor: props.editor,
        });

        if (props.clientRect === undefined || props.clientRect === null) {
          return;
        }

        const { clientRect } = props;

        const nextPopup = tippy('body', {
          getReferenceClientRect: clientRect as () => DOMRect,
          appendTo: () => document.body,
          content: component.element,
          showOnCreate: true,
          interactive: true,
          trigger: 'manual',
          placement: 'bottom-start',
        });

        if (nextPopup.length > 0 && nextPopup[0] !== undefined) {
          popup = nextPopup[0];
        }
      },

      onUpdate(props) {
        component.updateProps(props);

        if (!props.clientRect) {
          return;
        }

        popup.setProps({
          getReferenceClientRect: props.clientRect as () => DOMRect,
        });
      },

      onKeyDown(props) {
        if (props.event.key === 'Escape') {
          popup.hide();

          return true;
        }

        return component.ref?.onKeyDown(props);
      },

      onExit() {
        popup.destroy();
        component.destroy();
      },
    };
  },
};
