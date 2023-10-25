import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

import Image from "../../../components/Image/Image";
import "./SinglePost.css";

const SinglePost = () => {
	const [state, setState] = useState<{
		title: string;
		author: string;
		date: string;
		image: string;
		content: string;
	}>({
		title: "",
		author: "",
		date: "",
		image: "",
		content: "",
	});

	const params = useParams();

	useEffect(() => {
		const postId = params.postId;
		fetch(`${import.meta.env.VITE_API_BASE_URL}/feed/post/${postId}`)
			.then((res) => {
				if (res.status !== 200) {
					throw new Error("Failed to fetch status");
				}
				return res.json();
			})
			.then((resData) => {
				console.log(resData);
				setState((prev) => ({
					...prev,
					title: resData.post.title,
					author: resData.post.creator.name,
					date: new Date(resData.post.createdAt).toLocaleDateString(
						"en-US"
					),
					image: `${import.meta.env.VITE_API_BASE_URL}/${
						resData.post.imageUrl
					}`,
					content: resData.post.content,
				}));
			})
			.catch((err) => {
				console.log(err);
			});
	}, [params]);

	return (
		<section className="single-post">
			<h1>{state.title}</h1>
			<h2>
				Created by {state.author} on {state.date}
			</h2>
			<div className="single-post__image">
				<Image contain imageUrl={state.image} />
			</div>
			<p>{state.content}</p>
		</section>
	);
};

export default SinglePost;
