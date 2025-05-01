<template>
  <div class="flex justify-center">
    <div class="space-y-4">
      <!-- 本机 IP 字段 -->
      <div>
        <label for="local-ip" class="label-glass">本机 IP</label>
        <input type="text" id="local-ip" class="input-glass text-center py-2" disabled v-model="localIp" />
      </div>

      <!-- 选择文件按钮 -->
      <div>
        <label class="label-glass">选择文件</label>
        <button @click="pick" class="input-glass text-center py-2">{{ filePath }}</button>
      </div>

      <!-- 接收方 IP -->
      <div>
        <label for="receiver-ip" class="label-glass">接收方 IP</label>
        <input type="text" id="receiver-ip" class="input-glass" v-model="receiverIp" @input="checkReceiverStatus" />
      </div>

      <!-- 接收方状态 -->
      <div>
        <label class="label-glass">接收方状态</label>
        <div :class="receiverStatusClass" class="input-glass text-center py-2">
          {{ receiverStatus }}
        </div>
      </div>

      <!-- 发送文件按钮 -->
      <div>
        <button @click="sendFile"
        class=" w-full px-6 py-3 rounded-full text-white font-semibold text-lg shadow-md transition-all duration-300 
          backdrop-blur-md bg-white/10 border border-white/20
          hover:bg-white/20 hover:shadow-[0_0_10px_rgba(255,255,255,0.3)] hover:border-white/40">
        🚀 发送文件
        </button>
      </div>

      <dialog ref="progressSendModalRef" class="modal">
        <div data-theme="dark" class="modal-box">
          <h3 class="text-lg font-bold">正在发送文件中...</h3>
          <progress class="progress progress-accent w-86 h-3" :value="sendProgress" max="100"></progress>
          <div class="text-center mt-4">{{ sendProgress }}%</div>
        </div>
      </dialog>

      
    </div>
  </div>
</template>

<script setup>
import { ref,onMounted } from 'vue'
import { GetInnerIP,CheckIP,SendFileMeta,ChooseFile } from '../../wailsjs/go/backend/App'
import {EventsOn} from '../../wailsjs/runtime/runtime.js'

const localIp = ref('none')  // 示例，本机 IP
const filePath = ref('未选择文件')
const receiverIp = ref('')
const receiverStatus = ref('未准备好')
const receiverStatusClass = ref('badge-warning') // 默认显示为黄色
const sendProgress = ref(0) // 进度条
const progressSendModalRef = ref(null) // 发送端进度条dialog


onMounted(async () => {
  try{
    localIp.value = await GetInnerIP()
  } catch(error){
    console.error('获取本机 IP 失败:', error)
  }finally{
    // 监听文件发送进度
    EventsOn('fileSendProgress', (progress) => {
      if (sendProgress.value==0) {
        progressSendModalRef.value.showModal()
      }
      sendProgress.value = progress
      if (sendProgress.value >=100){
        progressSendModalRef.value.close()
        sendProgress.value = 0
      }
    })
  }

})


//验证是否为合法的ipv4地址
function isValidIPv4(ip) {
  const ipv4Regex = /^(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)$/
  return ipv4Regex.test(ip)
}

// 模拟接收方是否准备好
async function checkReceiverStatus() {
  if (receiverIp.value === ''){
    receiverStatus.value = '未准备好'
    receiverStatusClass.value = 'badge-warning'
    return
  }else {
    if (isValidIPv4(receiverIp.value)) {
    try {
      const isReady = await CheckIP(receiverIp.value)
      if (isReady) {
        receiverStatus.value = '接收方准备好'
        receiverStatusClass.value = 'badge-success' // 状态为绿色
      } else {
        receiverStatus.value = '接收方未响应'
        receiverStatusClass.value = 'badge-error' // 状态为红色
      }
    } catch (error) {
      console.error('检测接收方状态失败:', error)
      receiverStatus.value = '检测失败'
      receiverStatusClass.value = 'badge-error' // 状态为红色
    }
  } else {
    receiverStatus.value = 'IP 格式无效'
    receiverStatusClass.value = 'badge-error' // 状态为红色
  }
  }
}

// 处理文件选择
const pick = async () => {
  try {
    const selectedFilePath = await ChooseFile()
    if (selectedFilePath && selectedFilePath !== '') {
      filePath.value = selectedFilePath
      console.log('选择的文件路径:', filePath.value)
    } else {
      filePath.value = '未选择文件' // 清空文件路径
      console.log(filePath.value)
    }
  } catch (error) {
    console.error('选择文件失败:', error)
    filePath.value = '未选择文件' // 清空文件路径
  }
}

function sendFile() {
  if (!filePath.value || filePath.value === '') {
    console.error('请先选择文件')
    alert('请先选择文件')
    return
  }

  if (!receiverIp.value || receiverIp.value === '') {
    console.error('请先输入接收方 IP')
    alert('请先输入接收方 IP')
    return
  }

  // 发送文件元数据
  SendFileMeta(filePath.value, receiverIp.value)
    .then(() => {
      console.log('文件名发送成功')
    })
    .catch((error) => {
      console.error('发送文件名失败:', error)
    })
  

}

</script>



<style scoped>
/* 自定义接收方状态样式 */
.badge-warning {
  background-color: #fbbf24 !important; /* 黄色 */
}

.badge-success {
  background-color: #10b981 !important; /* 绿色 */
}

.badge-error {
  background-color: #ef4444 !important; /* 红色 */
}

.input-glass {
  appearance: none;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 9999px;
  padding: 10px;
  color: white;
  font-size: 20px;
  cursor: pointer;
  transition: all 0.3s ease-in-out;
  min-height: 50px; /* 设置最小高度 */
}

.label-glass {
  @apply block text-xl font-semibold bg-clip-text text-transparent bg-gradient-to-r from-blue-300 via-green-600 to-purple-500;
}

</style>
