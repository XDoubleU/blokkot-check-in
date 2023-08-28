/* eslint-disable sonarjs/no-duplicate-string */
import { getCheckInsToday, getLocation, getMyUser, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import mockRouter from "next-router-mock"
import LocationDetail from "pages/settings/locations/[id]"
import { screen, render, waitFor } from "test-utils"
import { DefaultLocation, defaultUserMock, noUserMock } from "mocks"
import moment from "moment"
import { FULL_FORMAT } from "api-wrapper/types/apiTypes"

// eslint-disable-next-line max-lines-per-function
describe("LocationDetail (page)", () => {
  it("Show detailed information of location as default user", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)
    mocked(getUser).mockImplementation(defaultUserMock)

    mocked(getLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: DefaultLocation
      })
    })

    const time = new Date().toISOString()
    mocked(getCheckInsToday).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: [
          {
            id: 1,
            capacity: 10,
            createdAt: time,
            locationId: "locationId",
            schoolName: "Andere"
          }
        ]
      })
    })

    await mockRouter.push("/settings/locations/locationId")

    render(<LocationDetail />)

    await screen.findByRole("heading", { name: "location" })
    await screen.findByText(moment.utc(time).format(FULL_FORMAT))
  })

  it("Show detailed information of location as default user, no check-ins", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)
    mocked(getUser).mockImplementation(defaultUserMock)

    mocked(getLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: DefaultLocation
      })
    })

    mocked(getCheckInsToday).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: []
      })
    })

    await mockRouter.push("/settings/locations/locationId")

    render(<LocationDetail />)

    await screen.findByRole("heading", { name: "location" })
    await screen.findByText("Nothing to see here.")
  })

  it("Can't fetch detailed information of location", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)
    mocked(getUser).mockImplementation(defaultUserMock)

    mocked(getLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: false
      })
    })

    await mockRouter.push("/settings/locations/notALocationId")

    render(<LocationDetail />)

    await screen.findByText("User has no location")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/locations/locationId")

    render(<LocationDetail />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe(
      "/signin?redirect_to=%2Fsettings%2Flocations"
    )
  })
})
