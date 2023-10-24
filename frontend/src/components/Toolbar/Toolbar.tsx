import React from "react";

import "./Toolbar.css";

const toolbar = (props: { children: React.ReactNode }) => (
	<div className="toolbar">{props.children}</div>
);

export default toolbar;
