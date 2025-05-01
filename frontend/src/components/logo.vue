<template>
    <div class="flex justify-center mt-6">
        <div class="relative z-10 w-16 h-16">
          <div class="absolute inset-3 flex items-center justify-center animate-pulse">
            <svg viewBox="0 0 100 100" fill="none" stroke="url(#grad1)" stroke-width="5"
                stroke-linejoin="round"
                class="w-full h-full drop-shadow-[0_0_8px_#00c472] animate-spin-slow">
              <defs>
                <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color:#00c472;stop-opacity:1" />
                  <stop offset="100%" style="stop-color:#22d3ee;stop-opacity:1" />
                </linearGradient>
              </defs>
              <polygon points="50,5 61,38 95,38 67,59 78,91 50,70 22,91 33,59 5,38 39,38" />
            </svg>
          </div>

        </div>
    </div>
    

</template>

<style>
@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin-slow {
  animation: spin-slow 6s linear infinite;
}
</style>