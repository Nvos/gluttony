import StarterKit from '@tiptap/starter-kit';
import { SolidEditorContent, useEditor } from '@vrite/tiptap-solid';

// In case when mention/suggestion will be needed
// https://github.com/vriteio/vrite/tree/main/apps/web/src/lib/editor/extensions/slash-menu
// https://github.com/vriteio/tiptap-solid/issues/1

export const Editor = () => {
  const editor = useEditor({
    extensions: [StarterKit],
  });

  return (
    <div
      id="editor-root"
      style={{
        width: '100%',
        height: '800px',
      }}
    >
      <SolidEditorContent editor={editor()} />
    </div>
  );
};
