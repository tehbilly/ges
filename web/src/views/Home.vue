<template>
  <n-space vertical>
    <n-space style="padding: 24px;">
      <n-card size="medium" title="Order Actions">
        <n-button-group vertical size="large">
          <n-button @click="onCreateOrder" :disabled="currentOrderID !== undefined">Create Order</n-button>
          <n-button @click="onAssignOrder" :disabled="currentOrderID === undefined">Assign Order</n-button>
          <n-button @click="onEmit('OrderLeft')" :disabled="currentOrderID === undefined">emit: OrderLeft</n-button>
          <n-button @click="onEmit('OrderArrived')" :disabled="currentOrderID === undefined">emit: OrderArrived
          </n-button>
          <n-button @click="onEmit('OrderDelivered')" :disabled="currentOrderID === undefined">emit: OrderDelivered
          </n-button>
          <n-button @click="onReset" type="error" ghost :disabled="currentOrderID === undefined">Reset</n-button>
        </n-button-group>
      </n-card>
      <n-code v-if="currentOrderID !== undefined" language="json" :code="JSON.stringify(currentOrder, null, '  ')" style="margin: 24px;"/>
    </n-space>

    <n-card title="Order Events" class="event-container">
      <n-collapse arrow-placement="right">
        <n-collapse-item v-for="event in events" :title="event.event">
          <n-code language="json" :code="JSON.stringify(event.order, null, '  ')"/>
        </n-collapse-item>
      </n-collapse>
    </n-card>
  </n-space>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import {
  NSpace,
  NCard,
  NCode,
  NButton, NButtonGroup,
  NCollapse, NCollapseItem,
} from 'naive-ui';

type Event = {
  event: string
  order: unknown
};

const currentOrderID = ref<string>();
const currentOrder = ref<any>({});

const events = ref<Array<Event>>([]);

const socket = new WebSocket('ws://localhost:8888/listen');

socket.addEventListener('open', event => {
  console.log('WebSocket connection opened:', event);
});

socket.addEventListener('message', event => {
  const e = JSON.parse(event.data) as Event;
  console.log('Message from server:', e);
  events.value.push(e);
  currentOrder.value = e.order;
});

socket.addEventListener('error', event => {
  console.error('WebSocket error:', event);
});

socket.addEventListener('close', event => {
  console.log('WebSocket closed:', event);
});

async function onCreateOrder() {
  const response = await fetch('http://localhost:8888/create', { method: 'POST' });
  const orderInfo = await response.json();
  currentOrderID.value = orderInfo.OrderID;
}

async function onAssignOrder() {
  const response = await fetch('http://localhost:8888/assign', {
    method: 'POST',
    body: JSON.stringify({
      OrderID: currentOrderID.value,
      AssignTo: 'wmcgann',
    }),
  });
  console.log('Assignment response:', response);
}

async function onEmit(event: string) {
  const response = await fetch('http://localhost:8888/emit', {
    method: 'POST',
    body: JSON.stringify({
      OrderID: currentOrderID.value,
      EventName: event,
    }),
  });
}

function onReset() {
  currentOrderID.value = undefined;
  currentOrder.value = {};
  events.value = [];
}
</script>

<style>
.event-container {
  margin: 24px;
}
</style>
