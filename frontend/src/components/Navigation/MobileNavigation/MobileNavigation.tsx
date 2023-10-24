import NavigationItems from "../NavigationItems/NavigationItems";
import "./MobileNavigation.css";

const mobileNavigation = (props: {
	open: boolean;
	mobile: boolean;
	onChooseItem: () => void;
	isAuth: boolean;
	onLogout: () => void;
}) => (
	<nav className={["mobile-nav", props.open ? "open" : ""].join(" ")}>
		<ul
			className={["mobile-nav__items", props.mobile ? "mobile" : ""].join(
				" "
			)}
		>
			<NavigationItems
				mobile
				onChoose={props.onChooseItem}
				isAuth={props.isAuth}
				onLogout={props.onLogout}
			/>
		</ul>
	</nav>
);

export default mobileNavigation;
