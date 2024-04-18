import { Extension } from '@tiptap/core';
import { Plugin, PluginKey } from '@tiptap/pm/state';

export const ImageDndExtension = Extension.create({
  name: 'dnd',

  addProseMirrorPlugins() {
    return [
      new Plugin({
        key: new PluginKey('imageDnd'),
        props: {
          handleDrop: (view, event, slice, moved) => {
            console.log('dnd');
            if (moved || event.dataTransfer === null || event.dataTransfer.files.length !== 1) {
              return false;
            }

            const file = event.dataTransfer.files[0];
            if (file === undefined) {
              return false;
            }

            const img = new Image();
            img.src = window.URL.createObjectURL(file);
            img.onload = () => {
              const data = new FormData();
              data.append('file', file);

              fetch('http://localhost:6001/storage/', { method: 'POST', body: data }).then(
                (response) => {
                  response.text().then((url) => {
                    console.log('got url', url);
                    const img = new Image();
                    img.src = 'http://localhost:6001' + url;
                    console.log('src', img.src);
                    img.onload = () => {
                      console.log('dispatching');
                      const { schema } = view.state;
                      const coordinates = view.posAtCoords({
                        left: event.clientX,
                        top: event.clientY,
                      });
                      const node = schema.nodes!.image!.create({
                        src: 'http://localhost:6001' + url,
                      });
                      const transaction = view.state.tr.insert(coordinates!.pos, node!);
                      return view.dispatch(transaction);
                    };
                  });
                },
              );
            };

            return true;
          },
        },
      }),
    ];
  },
});
