package adapay

const SDK_VERSION = "1.3.1"

const BASE_URL = "https://api.adapay.tech"
const PAGE_BASE_URL = "https://page.adapay.tech"

const PAYMENT_CREATE = "/v1/payments"
const PAYMENT_QUERY = "/v1/payments/{payment_id}"
const PAYMENT_LIST_QUERY = "/v1/payments/list"
const PAYMENT_CLOSE = "/v1/payments/{payment_id}/close"
const PAYMENT_CONFIRM = "/v1/payments/confirm"
const PAYMENT_QUERY_CONFIRM = "/v1/payments/confirm/{payment_confirm_id}"
const PAYMENT_QUERY_CONFIRM_LIST = "/v1/payments/confirm/list"
const PAYMENT_REVERSE = "/v1/payments/reverse"
const PAYMENT_QUERY_REVERSE = "/v1/payments/reverse/{reverse_id}"
const PAYMENT_QUERY_REVERSE_LIST = "/v1/payments/reverse/list"

const REFUND_CREATE = "/v1/payments/{payment_id}/refunds"
const REFUND_QUERY = "/v1/payments/refunds"

const BILL_DOWNLOAD = "/v1/bill/download"

const MEMBER_CREATE = "/v1/members"
const MEMBER_QUERY = "/v1/members/{member_id}"
const MEMBER_QUERY_LIST = "/v1/members/list"
const MEMBER_UPDATE = "/v1/members/update"

const SETTLE_ACCOUNT_CREATE = "/v1/settle_accounts"
const SETTLE_ACCOUNT_QUERY = "/v1/settle_accounts/{settle_account_id}"
const SETTLE_ACCOUNT_DELETE = "/v1/settle_accounts/delete"
const SETTLE_ACCOUNT_DETAIL_QUERY = "/v1/settle_accounts/settle_details"
const SETTLE_ACCOUNT_MODIFY = "/v1/settle_accounts/modify"
const SETTLE_ACCOUNT_BALANCE = "/v1/settle_accounts/balance"
const SETTLE_ACCOUNT_COMMISSIONS = "/v1/settle_accounts/commissions"
const SETTLE_ACCOUNT_COMMISSIONS_LIST = "/v1/settle_accounts/commissions/list"

const CORP_MEMBERS_CREATE = "/v1/corp_members"
const CORP_MEMBERS_QUERY = "/v1/corp_members/{member_id}"

const CREATE_CASHS = "/v1/cashs"
const QUERY_CASHS_STAT = "/v1/cashs/stat"

const USER_IDENTITY = "/v1/union/user_identity"

const WALLET_LOGIN = "/v1/walletLogin"
const WALLET_PAY = "/v1/account/payment"

const WALLET_CHECKOUT = "/v1/checkout"
const QUERY_CHECKOUT_LIST = "/v1/checkout/list"

const FAST_CARD_APPLY = "/v1/fast_card/apply"
const FAST_CARD_CONFIRM = "/v1/fast_card/confirm"
const FAST_CARD_LIST = "/v1/fast_card/list"

const FAST_PAY_CONFIRM = "/v1/fast_pay/confirm"
const FAST_PAY_SMS_CODE = "/v1/fast_pay/sms_code"

const FREEZE_ACCOUNT_FREEZE = "/v1/settle_accounts/freeze"
const FREEZE_ACCOUNT_FREEZE_LIST = "/v1/settle_accounts/freeze/list"

const UnFREEZE_ACCOUNT_FREEZE = "/v1/settle_accounts/unfreeze"
const UnFREEZE_ACCOUNT_FREEZE_LIST = "/v1/settle_accounts/unfreeze/list"

const TRANSFER_CREATE = "/v1/settle_accounts/transfer"
const TRANSFER_LIST = "/v1/settle_accounts/transfer/list"
