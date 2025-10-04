import apiRequest from "./client";

type Payment = {
  paymentId: number;
  leaseId: number;
  paymentAmount: number;
  paymentDateUnix: number;
  paymentNotes: string;
  paymentConfirmation: Uint8Array;
};

export async function createPayment(payment:Payment) {
  return apiRequest("/payments", "POST", payment);
}

export async function getPayments() {
  return apiRequest("/payments/", "GET");
}

export async function updatePayment(payment:Payment) {
  return apiRequest("/payments/update", "PUT", payment);
}

export async function deletePayment(paymentId:number) {
  return apiRequest(`/payments/delete/${paymentId}`, "DELETE");
}
