# User-service TODOs:
* Adicionar endpoint de atualização do refresh token
* Adicionar endpoint de logout
* Para o Refresh Token, seria ideal salvar o jti (JWT ID) em um Redis no user-service. Assim, se o usuário fizer logout, você apaga o ID do Redis, invalidando o Refresh Token imediatamente.