<template>
  <div class="relative min-h-screen bg-[url('/images/asfalt-light.png')] bg-[#004735]">
    <div class="container mx-auto p-6 relative z-10">

      <logo/>
      <send />
      
      <dialog ref="modalRef" class="modal">
        <div data-theme="dark" class="modal-box">
          <h3 class="text-lg font-bold">收到发送文件请求</h3>
          <input type="text" placeholder="Primary" class="input input-primary" v-model="fileName" />
          <div class="stat-value text-primary">{{ fileSize }}</div>
        
          <div class="modal-action">
            <button class="btn btn-success" @click="confirmReceive">接受</button>

            <form method="dialog">
              <!-- if there is a button in form, it will close the modal -->
              <button class="btn btn-error">拒绝</button>
            </form>

          </div>
        </div>
      </dialog>

      <dialog ref="progressReceiveModalRef" class="modal">
        <div data-theme="dark" class="modal-box">
          <h3 class="text-lg font-bold">正在接收文件中...</h3>
          <progress class="progress progress-accent w-86 h-3" :value="progress" max="100"></progress>
          <div class="text-center mt-4">{{ progress }}%</div>
        </div>
      </dialog>

      
    </div>

  </div>
</template>

<script setup>
import send from './components/send.vue'
import logo from './components/logo.vue'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { EventsOn,EventsOff } from '../wailsjs/runtime/runtime'
import { AcceptFileRequest } from '../wailsjs/go/backend/App'

const fileName = ref('')
const fileSize = ref('') // 修改为字符串类型
const fileLength = ref(0) // 文件大小(字节数)
const modalRef = ref(null) // 是否接收文件dialog
const progress = ref(0)  //进度条
const progressReceiveModalRef = ref(null) // 接收端进度条dialog


// 格式化文件大小
function formatFileSize(bytes) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

onMounted(() => {

  // 监听后端发送的文件请求
  EventsOn('fileInfo', (fileInfo) => {
    console.log('接收到文件信息:', fileInfo)
    fileName.value = fileInfo.fileName
    fileSize.value = formatFileSize(fileInfo.fileSize) // 转换文件大小为人类可读格式
    fileLength.value = fileInfo.fileSize // 保存原始文件大小
    modalRef.value.showModal()
  })

  //监听文件接收进度
  EventsOn('fileReceiveProgress', (progressValue) => {
  if (progress.value === 0) {
    progressReceiveModalRef.value.showModal() // 第一次接收到进度时弹窗
  }
  progress.value = progressValue

  if (progressValue >= 100) {
    progress.value = 0 // 重置进度条
    progressReceiveModalRef.value.close() // 完成后关闭进度弹窗
  }
})

  EventsOn('fileSaveSuccess',(status) => {
    console.log(status)
  })


})

onBeforeUnmount(() => {
  // 清理事件，防止内存泄漏
  EventsOff('fileInfo')
  EventsOff('fileReceiveProgress')
  EventsOff('fileSaveSuccess')
})

function confirmReceive() {
  
  modalRef.value.close()
  // 处理接收逻辑
  AcceptFileRequest(fileName.value,fileLength.value)
  
  
}
</script>


