import { Fragment, useEffect, useState, useCallback } from "react";

import Post from "../../components/Feed/Post/Post";
import Button from "../../components/Button/Button";
import FeedEdit from "../../components/Feed/FeedEdit/FeedEdit";
import Input from "../../components/Form/Input/Input";
import Paginator from "../../components/Paginator/Paginator";
import Loader from "../../components/Loader/Loader";
import ErrorHandler from "../../components/ErrorHandler/ErrorHandler";
import "./Feed.css";

interface post {
	_id: string;
	creator: {
		name: string;
	};
	createdAt: string;
	title: string;
	imageUrl?: string;
	content: string;
}

const Feed = () => {
	const [state, setState] = useState<{
		isEditing: boolean;
		posts: post[];
		totalPosts: number;
		editPost: post | null;
		status: string;
		postPage: number;
		postsLoading: boolean;
		editLoading: boolean;
		error: Error | null;
	}>({
		isEditing: false,
		posts: [],
		totalPosts: 0,
		editPost: null,
		status: "",
		postPage: 1,
		postsLoading: true,
		editLoading: false,
		error: null,
	});

	const loadPosts = useCallback(
		(direction?: "next" | "previous") => {
			if (direction) {
				setState((prev) => ({
					...prev,
					postsLoading: true,
					posts: [],
				}));
			}
			let page = state.postPage;
			if (direction === "next") {
				page++;
				setState((prev) => ({ ...prev, postPage: page }));
			}
			if (direction === "previous") {
				page--;
				setState((prev) => ({ ...prev, postPage: page }));
			}
			fetch(`${import.meta.env.VITE_API_BASE_URL}/feed/posts`)
				.then((res) => {
					if (res.status !== 200) {
						throw new Error("Failed to fetch posts.");
					}
					return res.json();
				})
				.then((resData) => {
					setState((prev) => ({
						...prev,
						posts: resData,
						totalPosts: resData.length,
						postsLoading: false,
					}));
				})
				.catch(catchError);
		},
		[state.postPage]
	);

	useEffect(() => {
		fetch("URL")
			.then((res) => {
				if (res.status !== 200) {
					throw new Error("Failed to fetch user status.");
				}
				return res.json();
			})
			.then((resData) => {
				setState((prev) => ({ ...prev, status: resData.status }));
			})
			.catch(catchError);

		loadPosts();
	}, [loadPosts]);

	const statusUpdateHandler = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		fetch("URL")
			.then((res) => {
				if (res.status !== 200 && res.status !== 201) {
					throw new Error("Can't update status!");
				}
				return res.json();
			})
			.then((resData) => {
				console.log(resData);
			})
			.catch(catchError);
	};

	const newPostHandler = () => {
		setState((prev) => ({ ...prev, isEditing: true }));
	};

	const startEditPostHandler = (postId: string) => {
		setState((prevState) => {
			const post = prevState.posts.find((p) => p._id === postId);
			if (!post) {
				return prevState; // or handle this situation differently
			}

			const loadedPost = { ...post };

			return {
				...prevState,
				isEditing: true,
				editPost: loadedPost,
			};
		});
	};

	const cancelEditHandler = () => {
		setState((prev) => ({ ...prev, isEditing: false, editPost: null }));
	};

	const finishEditHandler = (postData: post) => {
		setState((prev) => ({ ...prev, editLoading: true }));
		// Set up data (with image!)
		let url = "URL";
		if (state.editPost) {
			url = "URL";
		}

		fetch(url)
			.then((res) => {
				if (res.status !== 200 && res.status !== 201) {
					throw new Error("Creating or editing a post failed!");
				}
				return res.json();
			})
			.then((resData) => {
				const post = {
					_id: resData.post._id,
					title: resData.post.title,
					content: resData.post.content,
					creator: resData.post.creator,
					createdAt: resData.post.createdAt,
				};
				setState((prevState) => {
					let updatedPosts = [...prevState.posts];
					if (prevState.editPost) {
						const postIndex = prevState.posts.findIndex(
							(p) => p._id === prevState.editPost?._id
						);
						updatedPosts[postIndex] = post;
					} else if (prevState.posts.length < 2) {
						updatedPosts = prevState.posts.concat(post as post);
					}
					return {
						...prevState,
						posts: updatedPosts,
						isEditing: false,
						editPost: null,
						editLoading: false,
					};
				});
			})
			.catch((err) => {
				console.log(err);
				setState((prev) => ({
					...prev,
					isEditing: false,
					editPost: null,
					editLoading: false,
					error: err,
				}));
			});
	};

	const statusInputChangeHandler = (input: any, value: string) => {
		setState((prev) => ({ ...prev, status: value }));
	};

	const deletePostHandler = (postId: string) => {
		setState((prev) => ({ ...prev, postsLoading: true }));
		fetch("URL")
			.then((res) => {
				if (res.status !== 200 && res.status !== 201) {
					throw new Error("Deleting a post failed!");
				}
				return res.json();
			})
			.then((resData) => {
				console.log(resData);
				setState((prevState) => {
					const updatedPosts = prevState.posts.filter(
						(p) => p._id !== postId
					);
					return {
						...prevState,
						posts: updatedPosts,
						postsLoading: false,
					};
				});
			})
			.catch((err) => {
				console.log(err);
				setState((prev) => ({ ...prev, postsLoading: false }));
			});
	};

	const errorHandler = () => {
		setState((prev) => ({ ...prev, error: null }));
	};

	const catchError = (error: Error) => {
		setState((prev) => ({ ...prev, error: error }));
	};

	return (
		<Fragment>
			<ErrorHandler error={state.error} onHandle={errorHandler} />
			<FeedEdit
				editing={state.isEditing}
				selectedPost={state.editPost}
				loading={state.editLoading}
				onCancelEdit={cancelEditHandler}
				onFinishEdit={finishEditHandler}
			/>
			<section className="feed__status">
				<form onSubmit={statusUpdateHandler}>
					<Input
						type="text"
						placeholder="Your status"
						control="input"
						onChange={statusInputChangeHandler}
						value={state.status}
					/>
					<Button mode="flat" type="submit">
						Update
					</Button>
				</form>
			</section>
			<section className="feed__control">
				<Button mode="raised" design="accent" onClick={newPostHandler}>
					New Post
				</Button>
			</section>
			<section className="feed">
				{state.postsLoading && (
					<div style={{ textAlign: "center", marginTop: "2rem" }}>
						<Loader />
					</div>
				)}
				{state.posts.length <= 0 && !state.postsLoading ? (
					<p style={{ textAlign: "center" }}>No posts found.</p>
				) : null}
				{!state.postsLoading && (
					<Paginator
						onPrevious={loadPosts.bind(this, "previous")}
						onNext={loadPosts.bind(this, "next")}
						lastPage={Math.ceil(state.totalPosts / 2)}
						currentPage={state.postPage}
					>
						{state.posts.map((post) => (
							<Post
								key={post._id}
								id={post._id}
								author={post.creator.name}
								date={new Date(
									post.createdAt
								).toLocaleDateString("en-US")}
								title={post.title}
								image={post.imageUrl}
								content={post.content}
								onStartEdit={startEditPostHandler.bind(
									this,
									post._id
								)}
								onDelete={deletePostHandler.bind(
									this,
									post._id
								)}
							/>
						))}
					</Paginator>
				)}
			</section>
		</Fragment>
	);
};

export default Feed;