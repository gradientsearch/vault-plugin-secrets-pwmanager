class UserService {
	getEntityID() {
		let info = localStorage.getItem('loginInfo');
		if (info !== null) {
			let infoObj = JSON.parse(info);
            return  infoObj['entityID']
		}
	}
}


export const userService = new UserService()