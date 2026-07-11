import React from "react";
import {
ContextMenu,
ContextMenuContent,
ContextMenuItem,
ContextMenuLabel,
ContextMenuRadioGroup,
ContextMenuRadioItem,
ContextMenuSeparator,
ContextMenuShortcut,
ContextMenuSub,
ContextMenuSubContent,
ContextMenuSubTrigger,
ContextMenuTrigger
} from "~/components/ui/context-menu";

export default function FilesPageContext({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ContextMenu>
      <ContextMenuTrigger className="w-full h-full">
        {children}
      </ContextMenuTrigger>
      <ContextMenuContent className="w-64">
        <ContextMenuItem inset>
          Back
          <ContextMenuShortcut>⌘[</ContextMenuShortcut>
        </ContextMenuItem>
        <ContextMenuItem inset disabled>
          Forward
          <ContextMenuShortcut>⌘]</ContextMenuShortcut>
        </ContextMenuItem>
        <ContextMenuItem inset>
          Reload
          <ContextMenuShortcut>⌘R</ContextMenuShortcut>
        </ContextMenuItem>
        <ContextMenuItem inset>
          New File
          <ContextMenuShortcut>⌘]</ContextMenuShortcut>
        </ContextMenuItem>
        <ContextMenuItem inset>
          New Folder
          <ContextMenuShortcut>⌘]</ContextMenuShortcut>
        </ContextMenuItem>
        <ContextMenuSub>
          <ContextMenuSubTrigger inset>More Tools</ContextMenuSubTrigger>
          <ContextMenuSubContent className="w-48">
            <ContextMenuItem>
              Save Page As...
              <ContextMenuShortcut>⇧⌘S</ContextMenuShortcut>
            </ContextMenuItem>
            <ContextMenuItem>Create Shortcut...</ContextMenuItem>
            <ContextMenuItem>Name Window...</ContextMenuItem>
            <ContextMenuSeparator />
            <ContextMenuItem>Developer Tools</ContextMenuItem>
          </ContextMenuSubContent>
        </ContextMenuSub>

        <ContextMenuSeparator />
        <ContextMenuRadioGroup value="grid">
          <ContextMenuLabel inset>Layout</ContextMenuLabel>
          <ContextMenuSeparator />
          <ContextMenuRadioItem value="grid">Grid</ContextMenuRadioItem>
          <ContextMenuRadioItem value="list">List</ContextMenuRadioItem>
        </ContextMenuRadioGroup>
      </ContextMenuContent>
    </ContextMenu>
  );
}
