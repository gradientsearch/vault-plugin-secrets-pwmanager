class UserService {
	getEntityID() {
		let info = localStorage.getItem('loginInfo');
		if (info !== null) {
			let infoObj = JSON.parse(info);
            return  infoObj['entityID']
		}
	}

	getUsername() {
		let info = localStorage.getItem('loginInfo');
		if (info !== null) {
			let infoObj = JSON.parse(info);
            return  infoObj['username']
		}
	}
}


export const userService = new UserService()