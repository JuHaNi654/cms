import './style.css'
import htmx from "htmx.org"

declare global {
  interface Window {
    htmx: typeof htmx;
  }
}

window.htmx = htmx
//htmx.logAll()

type DragAreaFunction = (item: DragItem) => void

/* Drag and Drop scripts */
class DragItem {
  elem: Element
  toolbox: HTMLElement
  dragBtn?: HTMLButtonElement
  cbDragStart: DragAreaFunction
  cbDragEnd: DragAreaFunction
  cbMoveItem: DragAreaFunction
  dragging = false

  prevItem: DragItem | undefined
  nextItem: DragItem | undefined

  constructor(
    element: Element, 
    parentDragStart: DragAreaFunction, parentDragEnd: DragAreaFunction,
    moveItem: DragAreaFunction
  ) {
    this.elem = element
    this.toolbox = element.querySelector('[data-target="toolbox"]') as HTMLElement 
    this.cbDragStart = parentDragStart
    this.cbDragEnd = parentDragEnd
    this.cbMoveItem = moveItem
  }

  clone = (): DragItem => {
    const cloneNode = this.elem.cloneNode(true) as Element
    return new DragItem(
      cloneNode,
      this.cbDragStart, this.cbDragEnd, this.cbMoveItem
    )
  }

  setDragOn = () => {
    this.elem.setAttribute("draggable", "true")
    this.elem.classList.add("drag-active")
    this.dragging = true
 
    if (this.cbDragStart) {
      this.cbDragStart(this)
    }
  }

  matchElement = (element: Element): this | undefined =>  {
    if (element === this.elem) {
      return this
    }

    return undefined
  }

  setDragOff = () => {
    this.elem.setAttribute("draggable", "false")
    this.elem.classList.remove("drag-active") 
    this.dragging = false 
  
    if (this.cbDragEnd) {
      this.cbDragEnd(this)
    }
  } 

  // FIXME: Drag event bugs out when clicked from padding area
  onDragStart = (e: DragEvent) =>  {
    const x = this.dragBtn!.offsetWidth / 2 
    const y = this.dragBtn!.offsetHeight / 2
    e.dataTransfer?.setDragImage(this.dragBtn as HTMLButtonElement, x, y) 
  }

  onDragEnd = (_: DragEvent) => {
    this.setDragOff() 
  }

  // Prevent default, so that we can enable drop event on it
  onDragOver = (e: DragEvent) => {
    e.preventDefault()
  }

  onDragDrop = (e: DragEvent) => {
    e.preventDefault() // Prevent to open links in some elements
    const target = e.target as Element 

    if (target.classList.contains("drag-drop-zone")) {
      this.cbMoveItem(this) 
    }
  }

  setToolboxItems = () => {
    this.dragBtn = this.toolbox.querySelector('[data-target="toolbox-move"]') as HTMLButtonElement
    this.dragBtn.addEventListener("mousedown", this.setDragOn)
    this.dragBtn.addEventListener("mouseup", this.setDragOff)
  }

  init = (): this => {
    this.setToolboxItems()
    this.elem.addEventListener("dragstart", this.onDragStart as EventListener)
    this.elem.addEventListener("dragend", this.onDragEnd as EventListener)
    this.elem.addEventListener("dragover", this.onDragOver as EventListener)
    this.elem.addEventListener("drop", this.onDragDrop as EventListener)

    return this
  }
}

class DragArea {
  element: Element
  items: DragItem | undefined = undefined 
  activeItem: DragItem | undefined = undefined

  constructor(element: Element) {
    this.element = element
  }

  setActiveItem = (item: DragItem) => this.activeItem = item

  itemOnDrag = (item: DragItem) => {
    this.activeItem = item 
    let node = this.items
    while (node) {
      if (node !== item) {
        node.elem.classList.add("drag-drop-zone")
      }

      node = node.nextItem
    }
  }

  itemOnEnd = (_: DragItem) => {
    this.activeItem = undefined
    let node = this.items 
    while (node) {
      node.elem.classList.remove("drag-drop-zone")
      node = node.nextItem
    }
  }

  moveItem = (target: DragItem) => {
    if (!this.activeItem) return
    target.elem.before(this.activeItem.elem) 

    // Check if we are moving first item in the list 
    let node = this.items

    while (node) { 
      // Remove activeItem from list 
      if (this.activeItem === node) {
        if (this.activeItem === this.items) {
          const nextTmp = this.activeItem.nextItem
          nextTmp!.prevItem = undefined 
          this.items = nextTmp 
          this.activeItem.nextItem = undefined
        } else {
          const nextTmp = this.activeItem.nextItem
          const prevTmp = this.activeItem.prevItem
         
          if (nextTmp) {
            nextTmp.prevItem = prevTmp
          }
          prevTmp!.nextItem = nextTmp

          this.activeItem.prevItem = this.activeItem.nextItem = undefined
        } 
      }

      node = node.nextItem
    }

    node = this.items 

    while (node) {
      if (target === node) {
        if (target === this.items) {
          node.prevItem = this.activeItem
          this.activeItem.nextItem = node
          this.items = this.activeItem
        } else {
          const tmpPrev = node.prevItem
          node.prevItem = this.activeItem
          tmpPrev!.nextItem = this.activeItem
          this.activeItem.nextItem = node 
          this.activeItem.prevItem = tmpPrev
        }
      }
      node = node.nextItem
    }

  }

  logNodes = () => {
    let node = this.items 
    while (node) { 
      console.log("Node: ", node)
      node = node.nextItem
    }
  }

  init = (): this => {
    const items = document.querySelectorAll('.editor__item')
    let previous: DragItem | undefined = undefined 
    for (let i = 0; i < items.length; i++) {
      const dragItem = new DragItem(
        items[i],
        this.itemOnDrag,
        this.itemOnEnd,
        this.moveItem,
      ).init() 

      if (i === 0) {
        this.items = dragItem
        previous = dragItem
        continue
      }
 
      dragItem.prevItem = previous
      previous!.nextItem = dragItem
      previous = dragItem
    } 
    
    return this
  }
}

document.addEventListener("DOMContentLoaded", () => {
  const editor = document.querySelector(".editor-area")
  if (editor) new DragArea(editor).init()
})

