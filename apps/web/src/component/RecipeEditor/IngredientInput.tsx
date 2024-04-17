import { NodeViewWrapper, NodeViewProps, NodeViewContent } from '@tiptap/react';
import { css } from '~/ui/css';

export const IngredientInput = (props: NodeViewProps) => {
  console.log(props);

  return (
    <NodeViewWrapper className={css({ position: 'relative', padding: '400' })}>
      <NodeViewContent
        as="div"
        className={css({
          border: 'dotted 1px red',
        })}
      />
    </NodeViewWrapper>
  );
};
