import apiRequest from "./client";

// Payment record shape used by the API layer.
type Payment = {
  paymentId: number;
  leaseId: number;
  paymentAmount: number;
  paymentDateUnix: number; 
  paymentNotes: string;
  paymentConfirmation: Uint8Array; 
};

/**
 * Create a payment record.
 * @param payment - Payment payload
 */
export async function createPayment(payment: Payment) {
  return apiRequest("/payments", "POST", payment);
}

/**
 * Retrieve payments. The server may return all payments or a filtered set
 * depending on the endpoint implementation.
 */
export async function getPayments() {
  return apiRequest("/payments/", "GET");
}

/**
 * Update an existing payment record.
 * @param payment - Payment payload including `paymentId`
 */
export async function updatePayment(payment: Payment) {
  return apiRequest("/payments/update", "PUT", payment);
}

/**
 * Delete a payment by id.
 * @param paymentId - Numeric id of the payment to remove
 */
export async function deletePayment(paymentId: number) {
  return apiRequest(`/payments/delete/${paymentId}`, "DELETE");
}
