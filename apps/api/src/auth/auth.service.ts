import { Injectable } from "@nestjs/common"
import { UsersService } from "../users/users.service"
import { type Tokens } from "types-custom"
import { JwtService } from "@nestjs/jwt"
import { type Response } from "express"
import { compareSync } from "bcrypt"
import { type UserEntity } from "mikro-orm-config"

@Injectable()
export class AuthService {
  private readonly usersService: UsersService
  private readonly jwtService: JwtService

  public constructor(usersService: UsersService, jwtService: JwtService) {
    this.usersService = usersService
    this.jwtService = jwtService
  }

  public async signin(
    username: string,
    password: string
  ): Promise<Tokens | null> {
    const user = await this.usersService.getByUserName(username)
    if (!user) {
      return null
    }

    if (!compareSync(password, user.passwordHash)) {
      return null
    }

    return await this.getTokens(user)
  }

  public async refreshTokens(user: UserEntity): Promise<Tokens> {
    return this.getTokens(user)
  }

  public setTokensAsCookies(tokens: Tokens, res: Response): void {
    const accessTokenExpires = parseInt(
      (this.jwtService.decode(tokens.accessToken) as Record<string, string>).exp
    )
    const refreshTokenExpires = parseInt(
      (this.jwtService.decode(tokens.refreshToken) as Record<string, string>)
        .exp
    )

    res.cookie("accessToken", tokens.accessToken, {
      expires: new Date(accessTokenExpires * 1000),
      sameSite: "strict",
      httpOnly: true,
      secure: process.env.NODE_ENV === "production"
    })
    res.cookie("refreshToken", tokens.refreshToken, {
      expires: new Date(refreshTokenExpires * 1000),
      sameSite: "strict",
      httpOnly: true,
      path: "/auth/refresh",
      secure: process.env.NODE_ENV === "production"
    })
  }

  public async getTokens(user: UserEntity): Promise<Tokens> {
    const [accessToken, refreshToken] = await Promise.all([
      this.jwtService.signAsync(
        {
          sub: user.id
        },
        {
          secret: process.env.JWT_ACCESS_SECRET ?? "",
          expiresIn: process.env.JWT_ACCESS_EXPIRATION ?? ""
        }
      ),
      this.jwtService.signAsync(
        {
          sub: user.id
        },
        {
          secret: process.env.JWT_REFRESH_SECRET ?? "",
          expiresIn: process.env.JWT_REFRESH_EXPIRATION ?? ""
        }
      )
    ])

    return {
      accessToken,
      refreshToken
    }
  }
}
