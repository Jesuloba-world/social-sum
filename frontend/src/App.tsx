import { Fragment, useState, useEffect, useCallback } from "react";
import { Route, useNavigate, Routes, Navigate } from "react-router-dom";

import Layout from "./components/Layout/Layout";
import Backdrop from "./components/Backdrop/Backdrop";
import Toolbar from "./components/Toolbar/Toolbar";
import MainNavigation from "./components/Navigation/MainNavigation/MainNavigation.tsx";
import MobileNavigation from "./components/Navigation/MobileNavigation/MobileNavigation";
import ErrorHandler from "./components/ErrorHandler/ErrorHandler";
import FeedPage from "./pages/Feed/Feed";
import SinglePostPage from "./pages/Feed/SinglePost/SinglePost";
import LoginPage from "./pages/Auth/Login";
import SignupPage from "./pages/Auth/Signup";
import "./App.css";

const App = () => {
	const [showBackdrop, setShowBackdrop] = useState(false);
	const [showMobileNav, setShowMobileNav] = useState(false);
	const [isAuth, setIsAuth] = useState(false);
	const [token, setToken] = useState<string | null>(null);
	const [userId, setUserId] = useState<string | null>(null);
	const [authLoading, setAuthLoading] = useState(false);
	const [error, setError] = useState<Error | null>(null);

	const navigate = useNavigate();

	const setAutoLogout = useCallback((milliseconds: number) => {
		setTimeout(() => {
			logoutHandler();
		}, milliseconds);
	}, []);

	useEffect(() => {
		const token = localStorage.getItem("token");
		const expiryDate = localStorage.getItem("expiryDate");
		if (!token || !expiryDate) {
			return;
		}
		if (new Date(expiryDate) <= new Date()) {
			logoutHandler();
			return;
		}
		const userId = localStorage.getItem("userId");
		const remainingMilliseconds =
			new Date(expiryDate).getTime() - new Date().getTime();
		setIsAuth(true);
		setToken(token);
		setUserId(userId);

		setAutoLogout(remainingMilliseconds);
	}, [setAutoLogout]);

	const mobileNavHandler = (isOpen: boolean) => {
		setShowMobileNav(isOpen);
		setShowBackdrop(isOpen);
	};

	const backdropClickHandler = () => {
		setShowBackdrop(false);
		setShowMobileNav(false);
		setError(null);
	};

	const logoutHandler = () => {
		setIsAuth(false);
		setToken(null);
		localStorage.removeItem("token");
		localStorage.removeItem("expiryDate");
		localStorage.removeItem("userId");
	};

	const loginHandler = (
		event: Event,
		authData: { email: string; password: string }
	) => {
		event.preventDefault();
		setAuthLoading(true);
		console.log(authData);
		fetch(`${import.meta.env.VITE_API_BASE_URL}/auth/login`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: authData.email,
				password: authData.password,
			}),
			credentials: "include",
		})
			.then((res) => {
				if (res.status === 422) {
					throw new Error("Validation failed.");
				}
				if (res.status !== 200 && res.status !== 201) {
					console.log("Error!");
					throw new Error("Could not authenticate you!");
				}
				return res.json();
			})
			.then((resData) => {
				console.log(resData);
				setIsAuth(true);
				setToken(resData.token);
				setAuthLoading(false);
				setUserId(resData.userId);

				localStorage.setItem("token", resData.token);
				localStorage.setItem("userId", resData.userid);
				const remainingMilliseconds = 60 * 60 * 1000;
				const expiryDate = new Date(
					new Date().getTime() + remainingMilliseconds
				);
				localStorage.setItem("expiryDate", expiryDate.toISOString());
				setAutoLogout(remainingMilliseconds);
			})
			.catch((err) => {
				console.log(err);
				setIsAuth(false);
				setAuthLoading(false);
				setError(err);
			});
	};

	const signupHandler = (
		event: Event,
		authData: {
			signupForm: {
				email: { value: string };
				password: { value: string };
				name: { value: string };
			};
		}
	) => {
		event.preventDefault();
		console.log(authData);
		setAuthLoading(true);
		fetch(`${import.meta.env.VITE_API_BASE_URL}/auth/signup`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: authData.signupForm.email.value,
				password: authData.signupForm.password.value,
				name: authData.signupForm.name.value,
			}),
		})
			.then((res) => {
				if (res.status === 422) {
					throw new Error(
						"Validation failed. Make sure the email address isn't used yet!"
					);
				}
				if (res.status !== 200 && res.status !== 201) {
					console.log("Error!");
					throw new Error("Creating a user failed!");
				}
				return res.json();
			})
			.then((resData) => {
				console.log(resData);
				setIsAuth(false);
				setAuthLoading(false);
				navigate("/", { replace: true });
			})
			.catch((err) => {
				console.log(err);
				setIsAuth(false);
				setAuthLoading(false);
				setError(err);
			});
	};

	const errorHandler = () => {
		setError(null);
	};

	let routes = (
		<Routes>
			<Route
				path="/"
				element={
					<LoginPage onLogin={loginHandler} loading={authLoading} />
				}
			/>
			<Route
				path="/signup"
				element={
					<SignupPage
						onSignup={signupHandler}
						loading={authLoading}
					/>
				}
			/>
			<Route path="/*" element={<Navigate to="/" />} />
		</Routes>
	);
	if (isAuth) {
		routes = (
			<Routes>
				<Route
					path="/"
					element={<FeedPage userId={userId} token={token} />}
				/>
				<Route
					path="/:postId"
					element={<SinglePostPage userId={userId} token={token} />}
				/>
				<Route path="/*" element={<Navigate to="/" />} />
			</Routes>
		);
	}
	return (
		<Fragment>
			{showBackdrop && <Backdrop onClick={backdropClickHandler} />}
			<ErrorHandler error={error} onHandle={errorHandler} />
			<Layout
				header={
					<Toolbar>
						<MainNavigation
							onOpenMobileNav={() => mobileNavHandler(true)}
							onLogout={logoutHandler}
							isAuth={isAuth}
						/>
					</Toolbar>
				}
				mobileNav={
					<MobileNavigation
						open={showMobileNav}
						mobile
						onChooseItem={() => mobileNavHandler(false)}
						onLogout={logoutHandler}
						isAuth={isAuth}
					/>
				}
			/>
			{routes}
		</Fragment>
	);
};

export default App;
