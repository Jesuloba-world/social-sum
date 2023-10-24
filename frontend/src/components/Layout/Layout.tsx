import { Fragment, ReactNode } from "react";

import "./Layout.css";

const layout = (props: {
	header: ReactNode;
	children?: ReactNode;
	mobileNav: ReactNode;
}) => (
	<Fragment>
		<header className="main-header">{props.header}</header>
		{props.mobileNav}
		<main className="content">{props.children}</main>
	</Fragment>
);

export default layout;
