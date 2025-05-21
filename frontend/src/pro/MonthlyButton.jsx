import { onMount } from "solid-js";

function MonthlyButton(props) {
    let stripeButtonRef = null;

    onMount(() => {
        if (stripeButtonRef) {
            stripeButtonRef.setAttribute("buy-button-id", import.meta.env.VITE_STRIPE_BUY_BUTTON_MONTHLY);
            stripeButtonRef.setAttribute("publishable-key", import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);
            stripeButtonRef.setAttribute("client-reference-id", props.user);
        }
    });
    return (
        <stripe-buy-button ref={stripeButtonRef}></stripe-buy-button>
    )
}


export default MonthlyButton;