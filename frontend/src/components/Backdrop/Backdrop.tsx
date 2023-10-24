import ReactDOM from "react-dom";

import "./Backdrop.css";

const backdrop = (props: { onClick: () => void; open?: boolean }) => {
	const element = document.getElementById("backdrop-root");
	if (!element) {
		return null; // or handle this situation differently
	}

	return ReactDOM.createPortal(
		<div
			className={["backdrop", props.open ? "open" : ""].join(" ")}
			onClick={props.onClick}
		/>,
		element
	);
};

export default backdrop;
